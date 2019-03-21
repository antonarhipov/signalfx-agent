package client

import (
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/signalfx-agent/internal/utils"
)

const (
	TransportStatsGroup  = "transport"
	HTTPStatsGroup       = "http"
	IndicesStatsGroup    = "indices"
	JVMStatsGroup        = "jvm"
	ThreadpoolStatsGroup = "thread_pool"
	ProcessStatsGroup    = "process"
)

// Method to fetch node stats
func GetNodeStatsDatapoints(nodeStatsOutput *NodeStatsOutput, defaultDims map[string]string, selectedThreadPools map[string]bool, nodeStatsGroupEnhancedOption map[string]bool) []*datapoint.Datapoint {
	var out []*datapoint.Datapoint
	for _, nodeStats := range nodeStatsOutput.NodeStats {
		out = append(out, GetNodeStatsDatapointsHelper(nodeStats, defaultDims, selectedThreadPools, nodeStatsGroupEnhancedOption)...)
	}
	return out
}

func GetNodeStatsDatapointsHelper(nodeStats NodeStats, defaultDims map[string]string, selectedThreadPools map[string]bool, nodeStatsGroupEnhancedOption map[string]bool) []*datapoint.Datapoint {
	var dps []*datapoint.Datapoint
	dps = append(dps, nodeStats.Jvm.fetchJVMStats(nodeStatsGroupEnhancedOption[JVMStatsGroup], defaultDims)...)
	dps = append(dps, nodeStats.Process.fetchProcessStats(nodeStatsGroupEnhancedOption[ProcessStatsGroup], defaultDims)...)
	dps = append(dps, nodeStats.Transport.fetchTransportStats(nodeStatsGroupEnhancedOption[TransportStatsGroup], defaultDims)...)
	dps = append(dps, nodeStats.HTTP.fetchHTTPStats(nodeStatsGroupEnhancedOption[HTTPStatsGroup], defaultDims)...)
	dps = append(dps, fetchThreadPoolStats(nodeStatsGroupEnhancedOption[ThreadpoolStatsGroup], nodeStats.ThreadPool, defaultDims, selectedThreadPools)...)
	//dps = append(dps, nodeStats.Indices.fetchIndexStats(false, nil)...)
	return dps
}

func (jvm *Jvm) fetchJVMStats(enhanced bool, dims map[string]string) []*datapoint.Datapoint {
	var out []*datapoint.Datapoint

	if enhanced {
		out = append(out, []*datapoint.Datapoint{
			prepareGaugeHelper("gauge.jvm.threads.count", dims, jvm.JvmThreadsStats.Count),
			prepareGaugeHelper("gauge.jvm.threads.peak_count", dims, jvm.JvmThreadsStats.PeakCount),
			prepareGaugeHelper("gauge.jvm.mem.heap-used-percent", dims, jvm.JvmMemStats.HeapUsedPercent),
			prepareGaugeHelper("gauge.jvm.mem.heap_max_in_bytes", dims, jvm.JvmMemStats.HeapMaxInBytes),
			prepareGaugeHelper("gauge.jvm.mem.non-heap-committed", dims, jvm.JvmMemStats.NonHeapCommittedInBytes),
			prepareGaugeHelper("gauge.jvm.mem.non-heap-used", dims, jvm.JvmMemStats.NonHeapUsedInBytes),
			prepareGaugeHelper("gauge.jvm.mem.pools.young.max_in_bytes", dims, jvm.JvmMemStats.Pools.Young.MaxInBytes),
			prepareGaugeHelper("gauge.jvm.mem.pools.young.used_in_bytes", dims, jvm.JvmMemStats.Pools.Young.UsedInBytes),
			prepareGaugeHelper("gauge.jvm.mem.pools.young.peak_used_in_bytes", dims, jvm.JvmMemStats.Pools.Young.PeakUsedInBytes),
			prepareGaugeHelper("gauge.jvm.mem.pools.young.peak_max_in_bytes", dims, jvm.JvmMemStats.Pools.Young.PeakMaxInBytes),
			prepareGaugeHelper("gauge.jvm.mem.pools.old.max_in_bytes", dims, jvm.JvmMemStats.Pools.Old.MaxInBytes),
			prepareGaugeHelper("gauge.jvm.mem.pools.old.used_in_bytes", dims, jvm.JvmMemStats.Pools.Old.UsedInBytes),
			prepareGaugeHelper("gauge.jvm.mem.pools.old.peak_used_in_bytes", dims, jvm.JvmMemStats.Pools.Old.PeakUsedInBytes),
			prepareGaugeHelper("gauge.jvm.mem.pools.old.peak_max_in_bytes", dims, jvm.JvmMemStats.Pools.Old.PeakMaxInBytes),
			prepareGaugeHelper("gauge.jvm.mem.pools.survivor.max_in_bytes", dims, jvm.JvmMemStats.Pools.Survivor.MaxInBytes),
			prepareGaugeHelper("gauge.jvm.mem.pools.survivor.used_in_bytes", dims, jvm.JvmMemStats.Pools.Survivor.UsedInBytes),
			prepareGaugeHelper("gauge.jvm.mem.pools.survivor.peak_used_in_bytes", dims, jvm.JvmMemStats.Pools.Survivor.PeakUsedInBytes),
			prepareGaugeHelper("gauge.jvm.mem.pools.survivor.peak_max_in_bytes", dims, jvm.JvmMemStats.Pools.Survivor.PeakMaxInBytes),
			prepareGaugeHelper("gauge.jvm.mem.buffer_pools.mapped.count", dims, jvm.BufferPools.Mapped.Count),
			prepareGaugeHelper("gauge.jvm.mem.buffer_pools.mapped.used_in_bytes", dims, jvm.BufferPools.Mapped.UsedInBytes),
			prepareGaugeHelper("gauge.jvm.mem.buffer_pools.mapped.total_capacity_in_bytes", dims, jvm.BufferPools.Mapped.TotalCapacityInBytes),
			prepareGaugeHelper("gauge.jvm.mem.buffer_pools.direct.count", dims, jvm.BufferPools.Direct.Count),
			prepareGaugeHelper("gauge.jvm.mem.buffer_pools.direct.used_in_bytes", dims, jvm.BufferPools.Direct.UsedInBytes),
			prepareGaugeHelper("gauge.jvm.mem.buffer_pools.direct.total_capacity_in_bytes", dims, jvm.BufferPools.Direct.TotalCapacityInBytes),
			prepareGaugeHelper("gauge.jvm.classes.current-loaded-count", dims, jvm.Classes.CurrentLoadedCount),

			prepareCumulativeHelper("counter.jvm.gc.count", dims, jvm.JvmGcStats.Collectors.Young.CollectionCount),
			prepareCumulativeHelper("counter.jvm.gc.old-count", dims, jvm.JvmGcStats.Collectors.Old.CollectionCount),
			prepareCumulativeHelper("counter.jvm.gc.old-time", dims, jvm.JvmGcStats.Collectors.Old.CollectionTimeInMillis),
			prepareCumulativeHelper("counter.jvm.classes.total-loaded-count", dims, jvm.Classes.TotalLoadedCount),
			prepareCumulativeHelper("counter.jvm.classes.total-unloaded-count", dims, jvm.Classes.TotalUnloadedCount),
		}...)
	}

	out = append(out, []*datapoint.Datapoint{
		prepareGaugeHelper("gauge.jvm.mem.heap-used", dims, jvm.JvmMemStats.HeapUsedInBytes),
		prepareGaugeHelper("gauge.jvm.mem.heap-committed", dims, jvm.JvmMemStats.HeapCommittedInBytes),
		prepareCumulativeHelper("counter.jvm.uptime", dims, jvm.UptimeInMillis),
		prepareCumulativeHelper("counter.jvm.gc.time", dims, jvm.JvmGcStats.Collectors.Young.CollectionTimeInMillis),
	}...)

	return out
}

func (processStats *Process) fetchProcessStats(enhanced bool, dims map[string]string) []*datapoint.Datapoint {
	var out []*datapoint.Datapoint

	if enhanced {
		out = append(out, []*datapoint.Datapoint{
			prepareGaugeHelper("counter.process.max_file_descriptors", dims, processStats.MaxFileDescriptors),
			prepareGaugeHelper("gauge.process.cpu.percent", dims, processStats.CPU.Percent),
			prepareCumulativeHelper("counter.process.cpu.time", dims, processStats.CPU.TotalInMillis),
			prepareCumulativeHelper("counter.process.mem.total-virtual-size", dims, processStats.Mem.TotalVirtualInBytes),
		}...)
	}

	out = append(out, []*datapoint.Datapoint{
		prepareGaugeHelper("gauge.process.open_file_descriptors", dims, processStats.OpenFileDescriptors),
	}...)

	return out
}

func fetchThreadPoolStats(enhanced bool, threadPools map[string]ThreadPoolStats, defaultDims map[string]string, selectedThreadPools map[string]bool) []*datapoint.Datapoint {
	var out []*datapoint.Datapoint
	for threadPool, stats := range threadPools {
		if !selectedThreadPools[threadPool] {
			continue
		}
		out = append(out, threadPoolDatapoints(enhanced, threadPool, stats, defaultDims)...)
	}
	return out
}

func threadPoolDatapoints(enhanced bool, threadPool string, threadPoolStats ThreadPoolStats, defaultDims map[string]string) []*datapoint.Datapoint {
	var out []*datapoint.Datapoint
	threadPoolDimension := map[string]string{}
	threadPoolDimension["thread_pool"] = threadPool

	dims := utils.MergeStringMaps(defaultDims, threadPoolDimension)

	if enhanced {
		out = append(out, []*datapoint.Datapoint{
			prepareCumulativeHelper("counter.thread_pool.rejected", dims, threadPoolStats.Threads),
			prepareCumulativeHelper("counter.thread_pool.rejected", dims, threadPoolStats.Queue),
			prepareCumulativeHelper("counter.thread_pool.rejected", dims, threadPoolStats.Active),
			prepareCumulativeHelper("counter.thread_pool.rejected", dims, threadPoolStats.Completed),
		}...)
	}

	out = append(out, []*datapoint.Datapoint{
		prepareCumulativeHelper("counter.thread_pool.rejected", dims, threadPoolStats.Rejected),
	}...)
	return out
}

func (transport *Transport) fetchTransportStats(enhanced bool, dims map[string]string) []*datapoint.Datapoint {
	var out []*datapoint.Datapoint

	if enhanced {
		out = append(out, []*datapoint.Datapoint{
			prepareGaugeHelper("gauge.transport.server_open", dims, transport.ServerOpen),
			prepareCumulativeHelper("counter.transport.rx.count", dims, transport.RxCount),
			prepareCumulativeHelper("counter.transport.rx.size", dims, transport.RxSizeInBytes),
			prepareCumulativeHelper("counter.transport.tx.count", dims, transport.TxCount),
			prepareCumulativeHelper("counter.transport.tx.size", dims, transport.TxSizeInBytes),
		}...)
	}
	return out
}

func (http *HTTP) fetchHTTPStats(enhanced bool, dims map[string]string) []*datapoint.Datapoint {
	var out []*datapoint.Datapoint

	if enhanced {
		out = append(out, []*datapoint.Datapoint{
			prepareGaugeHelper("gauge.http.current_open", dims, http.CurrentOpen),
			prepareCumulativeHelper("counter.http.total_open", dims, http.TotalOpened),
		}...)
	}

	return out
}
