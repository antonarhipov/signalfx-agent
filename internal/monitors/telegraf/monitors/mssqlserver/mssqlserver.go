package mssqlserver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	telegrafInputs "github.com/influxdata/telegraf/plugins/inputs"
	telegrafPlugin "github.com/influxdata/telegraf/plugins/inputs/sqlserver"
	"github.com/signalfx/signalfx-agent/internal/core/config"
	"github.com/signalfx/signalfx-agent/internal/monitors"
	"github.com/signalfx/signalfx-agent/internal/monitors/telegraf/common/accumulator"
	"github.com/signalfx/signalfx-agent/internal/monitors/telegraf/common/emitter"
	"github.com/signalfx/signalfx-agent/internal/monitors/telegraf/common/emitter/baseemitter"
	"github.com/signalfx/signalfx-agent/internal/monitors/telegraf/monitors/winperfcounters"
	"github.com/signalfx/signalfx-agent/internal/monitors/types"
	"github.com/signalfx/signalfx-agent/internal/utils"
	log "github.com/sirupsen/logrus"
)

const monitorType = "telegraf/sqlserver"

// MONITOR(telegraf/sqlserver): This monitor reports metrics about Microsoft SQL servers.
// This monitor is based on the telegraf sqlserver plugin.  More information about the telegraf plugin
// can be found [here](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/sqlserver).
//
// You will need to create a login on the SQL server for the monitor to use.  You can create this login by
// executing the following commands in an a SQL client while logged in as an administrator.
//
// ```
// USE master;
// GO
// CREATE LOGIN [signalfxagent] WITH PASSWORD = N'<YOUR PASSWORD HERE>';
// GO
// GRANT VIEW SERVER STATE TO [signalfxagent];
// GO
// GRANT VIEW ANY DEFINITION TO [signalfxagent];
// GO
// ```
//
// Troubleshooting:
//
// On some Windows based SQL server distributions TCP/IP has been disabled by default.  This behavior
// has been observed on Azure SQL server instances.  You may need to explicitly turn on TCP/IP for the
// SQL server if you see error messages simillar to the following.
//
// ```
// Cannot read handshake packet: read tcp: wsarecv: An existing connection was forcibly closed by the remote host.
// ```
//
// 1. Verify agent configurations are correct.
// 2. Ensure TCP/IP is enabled for the SQL server by going to `Start` -> `Administrative Tools` -> `Computer Management`
// 3. In the `Computer Management` side bar, drill down to `Services and Applications` -> `SQL Server Configuration Manager` -> `SQL Server Network Configuration`
// 4. Select `Protocols for <YOUR SQL SERVER NAME>`.
// 5. In the protocol list to the right, right-click on the `TCP/IP` protocol and `enable` it.
//
// Sample YAML configuration:
//
// ```yaml
// monitors:
//  - type: telegraf/sqlserver
//    host: hostname
//    port: 1433
//    userID: sa
//    password: P@ssw0rd!
//    appName: signalfxagent
// ```
//

var logger = log.WithFields(log.Fields{"monitorType": monitorType})

func init() {
	monitors.Register(monitorType, func() interface{} { return &Monitor{} }, &Config{})
}

// Config for this monitor
type Config struct {
	config.MonitorConfig `acceptsEndpoints:"true"`
	Host                 string `yaml:"host" validate:"required" default:"."`
	Port                 uint16 `yaml:"port" validate:"required" default:"1433"`
	// UserID used to access the SQL Server instance.
	UserID string `yaml:"userID"`
	// Password used to access the SQL Server instance.
	Password string `yaml:"password" neverLog:"true"`
	// The app name used by the monitor when connecting to the SQLServer.
	AppName string `yaml:"appName" default:"signalfxagent"`
	// The version of queries to use when accessing the cluster.
	// Please refer to the telegraf documentation for more information.
	QueryVersion int `yaml:"queryVersion" default:"2"`
	// Whether the database is an azure database or not.
	AzureDB bool `yaml:"azureDB"`
	// Queries to exclude possible values are `PerformanceCounters`, `WaitStatsCategorized`,
	// `DatabaseIO`, `DatabaseProperties`, `CPUHistory`, `DatabaseSize`, `DatabaseStats`, `MemoryClerk`
	// `VolumeSpace`, and `PerformanceMetrics`.
	ExcludeQuery []string `yaml:"excludedQueries"`
	// Log level to use when accessing the database
	Log uint `yaml:"log" default:"1"`
}

// Monitor for Utilization
type Monitor struct {
	Output types.Output
	cancel func()
}

// fetch the factory used to generate the perf counter plugin
var factory = telegrafInputs.Inputs["sqlserver"]

// Configure the monitor and kick off metric syncing
func (m *Monitor) Configure(conf *Config) error {
	plugin := factory().(*telegrafPlugin.SQLServer)

	server := fmt.Sprintf("Server=%s;Port=%d;", conf.Host, conf.Port)

	if conf.UserID != "" {
		server = fmt.Sprintf("%sUser Id=%s;", server, conf.UserID)
	}
	if conf.Password != "" {
		server = fmt.Sprintf("%sPassword=%s;", server, conf.Password)
	}
	if conf.AppName != "" {
		server = fmt.Sprintf("%sapp name=%s;", server, conf.AppName)
	}
	server = fmt.Sprintf("%slog=%d;", server, conf.Log)

	plugin.Servers = []string{server}
	plugin.QueryVersion = conf.QueryVersion
	plugin.AzureDB = conf.AzureDB
	plugin.ExcludeQuery = conf.ExcludeQuery

	// create batch emitter
	emit := baseemitter.NewEmitter(m.Output, logger)

	// Hard code the plugin name because the emitter will parse out the
	// configured measurement name as plugin and that is confusing.
	emit.AddTag("plugin", strings.Replace(monitorType, "/", "-", -1))

	// replacer sanitizes metrics according to our PCR reporter rules
	replacer := winperfcounters.NewPCRReplacer()

	emit.AddMeasurementTransformation(
		func(ms telegraf.Metric) error {
			// if it's a sqlserver_performance metric
			// remap the counter and value to a field
			if ms.Name() == "sqlserver_performance" {
				emitter.RenameFieldWithTag(ms, "counter", "value", replacer)
			}

			// if it's a sqlserver_memory_clerks metric remap clerk type to field
			if ms.Name() == "sqlserver_memory_clerks" {
				ms.SetName(fmt.Sprintf("sqlserver_memory_clerks.size_kb"))
				emitter.RenameFieldWithTag(ms, "clerk_type", "size_kb", replacer)
			}
			return nil
		})

	// convert the metric name to lower case
	emit.AddMetricNameTransformation(func(m string) string {
		return strings.ToLower(m)
	})

	// create the accumulator
	ac := accumulator.NewAccumulator(emit)

	// create contexts for managing the the plugin loop
	var ctx context.Context
	ctx, m.cancel = context.WithCancel(context.Background())

	// gather metrics on the specified interval
	utils.RunOnInterval(ctx, func() {
		if err := plugin.Gather(ac); err != nil {
			logger.WithError(err).Errorf("an error occurred while gathering metrics from the plugin")
		}

	}, time.Duration(conf.IntervalSeconds)*time.Second)

	return nil
}

// Shutdown stops the metric sync
func (m *Monitor) Shutdown() {
	if m.cancel != nil {
		m.cancel()
	}
}
