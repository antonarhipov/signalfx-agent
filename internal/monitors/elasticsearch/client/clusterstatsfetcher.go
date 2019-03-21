package client

import (
	"github.com/signalfx/golib/datapoint"
)

const (
	clusterStatusGreen = "green"
	clusterStatusYellow = "yellow"
	clsuterStatusRed = "red"
)

// GetClusterStatsDatapoints fetches datapoints for ES cluster level stats
func GetClusterStatsDatapoints(clusterStatsOutput *ClusterStatsOutput, defaultDims map[string]string, enhanced bool) []*datapoint.Datapoint {
	var out []*datapoint.Datapoint

	if enhanced {
		out = append(out, []*datapoint.Datapoint{
			prepareGaugeHelper("gauge.cluster.initializing-shards", defaultDims, clusterStatsOutput.InitializingShards),
			prepareGaugeHelper("gauge.cluster.delayed-unassigned-shards", defaultDims, clusterStatsOutput.DelayedUnassignedShards),
			prepareGaugeHelper("gauge.cluster.pending-tasks", defaultDims, clusterStatsOutput.NumberOfPendingTasks),
			prepareGaugeHelper("gauge.cluster.in-flight-fetches", defaultDims, clusterStatsOutput.NumberOfInFlightFetch),
			prepareGaugeHelper("gauge.cluster.task-max-wait-time", defaultDims, clusterStatsOutput.TaskMaxWaitingInQueueMillis),
			prepareGaugeFHelper("gauge.cluster.active-shards-percent", defaultDims, clusterStatsOutput.ActiveShardsPercentAsNumber),
			prepareGaugeHelper("gauge.cluster.status", defaultDims, getMetricValueFromClusterStatus(clusterStatsOutput.Status)),
		}...)
	}
	out = append(out, []*datapoint.Datapoint{
		prepareGaugeHelper("gauge.cluster.active-primary-shards", defaultDims, clusterStatsOutput.ActivePrimaryShards),
		prepareGaugeHelper("gauge.cluster.active-shards", defaultDims, clusterStatsOutput.ActiveShards),
		prepareGaugeHelper("gauge.cluster.number-of-data_nodes", defaultDims, clusterStatsOutput.NumberOfDataNodes),
		prepareGaugeHelper("gauge.cluster.number-of-nodes", defaultDims, clusterStatsOutput.NumberOfNodes),
		prepareGaugeHelper("gauge.cluster.relocating-shards", defaultDims, clusterStatsOutput.RelocatingShards),
		prepareGaugeHelper("gauge.cluster.unassigned-shards", defaultDims, clusterStatsOutput.UnassignedShards),
	}...)
	return out
}

// Map cluster status to a numeric value
func getMetricValueFromClusterStatus(s *string) *int64 {
	out := new(int64)
	status := *s

	switch status {
	case clusterStatusGreen:
		*out = 0
		break
	case clusterStatusYellow:
		*out = 1
		break
	case clsuterStatusRed:
		*out = 2
		break
	default:
		break
	}

	return out
}


