package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"github.com/signalfx/signalfx-agent/internal/core/config"
	"github.com/signalfx/signalfx-agent/internal/monitors"
	"github.com/signalfx/signalfx-agent/internal/monitors/elasticsearch/client"
	"github.com/signalfx/signalfx-agent/internal/monitors/types"
	"github.com/signalfx/signalfx-agent/internal/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

const monitorType = "elasticsearch"

var logger = log.WithFields(log.Fields{"monitorType": monitorType})

// Config for this monitor
type Config struct {
	config.MonitorConfig `yaml:",inline" acceptsEndpoints:"true"`
	Host                 string `yaml:"host" validate:"required"`
	Port                 string `yaml:"port" validate:"required"`
	// Username used to access elasticsearch stats api
	Username string `yaml:"username"`
	// Password used to access elasticsearch stats api
	Password string `yaml:"password" neverLog:"true"`
	// Protocol used to connect: http or https
	Protocol string `yaml:"protocol" default:"http"`
	// Cluster name to which the node belongs. This is an optional config that
	// will override the cluster name fetched from a node and will be used to
	// populate the plugin_instance dimension
	Cluster string `yaml:"cluster"`
	// Version of the elasticsearch cluster
	Version string `yaml:"version"`
	// EnableClusterHealth enables reporting on the cluster health
	EnableClusterHealth *bool `yaml:"enableClusterHealth" default:"true"`
	// Node stat groups to collect enhanced metrics from. By default the monitor
	// collects a subset of stats (see documentation for details) across the jvm,
	// process, threadpool, transport, http indices groups. Only valid options for
	// this option are transport, http, process, jvm, indices, thread_pool
	NodeStatsGroups []string `yaml:"enhancedNodeStats" default:"[]"`
	// ThreadPools to report threadpool node stats on
	ThreadPools []string `yaml:"threadPools"`
}

// Monitor for conviva metrics
type Monitor struct {
	Output  types.Output
	cancel  context.CancelFunc
	ctx     context.Context
	timeout time.Duration
}

func init() {
	monitors.Register(monitorType, func() interface{} { return &Monitor{} }, &Config{})
}

// Configure monitor
func (m *Monitor) Configure(conf *Config) error {
	fmt.Println(conf)
	esClient := client.NewESClient(conf.Host, conf.Port, conf.Protocol, conf.Username, conf.Password)
	m.ctx, m.cancel = context.WithCancel(context.Background())

	threadPools := []string{"search", "index"}
	threadPools = append(threadPools, conf.ThreadPools...)

	utils.RunOnInterval(m.ctx, func() {
		var nodeStatsOutput client.NodeStatsOutput
		err := esClient.GetNodeAndThreadPoolStats(&nodeStatsOutput)

		if err != nil {
			logger.Errorf("Failed to GET node stats", err)
			return
		}

		pluginInstanceDimension, err := prepareClusterDimension(conf.Cluster, nodeStatsOutput.ClusterName)

		if err != nil {
			logger.Errorf("Failed to create plugin_instance dimension", err)
			return
		}

		dps := client.GetNodeStatsDatapoints(&nodeStatsOutput, pluginInstanceDimension, threadPools, conf.NodeStatsGroups)

		for i := range dps {
			if dps[i] == nil {
				continue
			}
			m.Output.SendDatapoint(dps[i])
		}
	}, time.Duration(conf.IntervalSeconds)*time.Second)
	return nil
}

// Prepares dimensions that are common to all datapoints from the monitor
func prepareClusterDimension(userProvidedClusterName string, queriedClusterName *string) (map[string]string, error) {
	dims := map[string]string{}
	clusterName := userProvidedClusterName

	if clusterName == "" {
		if queriedClusterName == nil {
			return nil, errors.New("Failed to get cluster name from node")
		}
		clusterName = *queriedClusterName
	}

	dims["plugin_instance"] = clusterName

	return dims, nil
}

// Shutdown stops the metric sync
func (m *Monitor) Shutdown() {
	if m.cancel != nil {
		m.cancel()
	}
}
