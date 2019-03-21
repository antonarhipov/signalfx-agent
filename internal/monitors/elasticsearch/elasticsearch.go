package elasticsearch

import (
	"context"
	"errors"
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
	// Username used to access Elasticsearch stats API
	Username string `yaml:"username"`
	// Password used to access Elasticsearch stats API
	Password string `yaml:"password" neverLog:"true"`
	// Whether to use https or not
	UseHTTPS bool `yaml:"useHTTPS"`
	// Cluster name to which the node belongs. This is an optional config that
	// will override the cluster name fetched from a node and will be used to
	// populate the plugin_instance dimension
	Cluster string `yaml:"cluster"`
	// Version of the Elasticsearch cluster
	Version string `yaml:"version"`
	// EnableClusterHealth enables reporting on the cluster health
	EnableClusterHealth *bool `yaml:"enableClusterHealth" default:"true"`
	// Enable enhanced HTTP stats
	EnableEnhancedHTTPStats bool `yaml:"enableEnhancedHTTPStats"`
	// Enable enhanced Indices stats
	EnableEnhancedIndicesStats bool `yaml:"enableEnhancedIndicesStats"`
	// Enable enhanced JVM stats
	EnableEnhancedJVMStats bool `yaml:"enableEnhancedJVMStats"`
	// Enable enhanced Process stats
	EnableEnhancedProcessStats bool `yaml:"enableEnhancedProcessStats"`
	// Enable enhanced ThreadPool stats
	EnableEnhancedThreadPoolStats bool `yaml:"enableEnhancedThreadPoolStats"`
	// Enable enhanced Transport stats
	EnableEnhancedTransportStats bool `yaml:"enableEnhancedTransportStats"`
	// ThreadPools to report threadpool node stats on
	ThreadPools []string `yaml:"threadPools" default:"[\"search\", \"index\"]"`
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
	esClient := client.NewESClient(conf.Host, conf.Port, conf.UseHTTPS, conf.Username, conf.Password)
	m.ctx, m.cancel = context.WithCancel(context.Background())

	utils.RunOnInterval(m.ctx, func() {
		nodeStatsOutput, err, url := esClient.GetNodeAndThreadPoolStats()

		if err != nil {
			logger.WithError(err).Errorf("Failed to GET node stats. URL: %s", url)
			return
		}

		pluginInstanceDimension, err := prepareClusterDimension(conf.Cluster, nodeStatsOutput.ClusterName)

		if err != nil {
			logger.WithError(err).Errorf("Failed to prepare plugin_instance dimension")
			return
		}

		dps := client.GetNodeStatsDatapoints(nodeStatsOutput, pluginInstanceDimension, utils.StringSliceToMap(conf.ThreadPools), map[string]bool{
			client.HTTPStatsGroup:       conf.EnableEnhancedHTTPStats,
			client.IndicesStatsGroup:    conf.EnableEnhancedIndicesStats,
			client.JVMStatsGroup:        conf.EnableEnhancedJVMStats,
			client.ProcessStatsGroup:    conf.EnableEnhancedProcessStats,
			client.ThreadpoolStatsGroup: conf.EnableEnhancedThreadPoolStats,
			client.TransportStatsGroup:  conf.EnableEnhancedTransportStats,
		})

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
			return nil, errors.New("Failed to get cluster name from Elasticsearch API")
		}
		clusterName = *queriedClusterName
	}

	// "plugin_instance" dimension is added to maintain backwards compatibility with built-in content
	dims["plugin_instance"] = clusterName
	dims["cluster"] = clusterName

	return dims, nil
}

// Shutdown stops the metric sync
func (m *Monitor) Shutdown() {
	if m.cancel != nil {
		m.cancel()
	}
}
