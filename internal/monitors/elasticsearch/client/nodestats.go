package client

/*
 This nested struct contains all the node level stats of the following types -
 indices, process, jvm, thread pool, transport and http, and will be used to
 retrieve node level metrics
*/
type NodeStatsOutput struct {
	Nodes       Nodes                `json:"_nodes"`
	ClusterName *string              `json:"cluster_name"`
	NodeStats   map[string]NodeStats `json:"nodes"`
}

type Nodes struct {
	Total      *int64 `json:"total"`
	Successful *int64 `json:"successful"`
	Failed     *int64 `json:"failed"`
}

type NodeStats struct {
	Timestamp        *int64                     `json:"timestamp"`
	Name             *string                    `json:"name"`
	TransportAddress *string                    `json:"transport_address"`
	Host             *string                    `json:"host"`
	IP               *string                    `json:"ip"`
	Roles            *[]string                  `json:"roles"`
	Attributes       Attributes                 `json:"attributes"`
	Indices          IndexStatsGroups           `json:"indices"`
	Process          Process                    `json:"process"`
	Jvm              Jvm                        `json:"jvm"`
	ThreadPool       map[string]ThreadPoolStats `json:"thread_pool"`
	Transport        Transport                  `json:"transport"`
	HTTP             HTTP                       `json:"http"`
}

type ThreadPoolStats struct {
	Threads   *int64 `json:"threads"`
	Queue     *int64 `json:"queue"`
	Active    *int64 `json:"active"`
	Rejected  *int64 `json:"rejected"`
	Largest   *int64 `json:"largest"`
	Completed *int64 `json:"completed"`
}

type Process struct {
	Timestamp           *int64          `json:"timestamp"`
	OpenFileDescriptors *int64          `json:"open_file_descriptors"`
	MaxFileDescriptors  *int64          `json:"max_file_descriptors"`
	CPU                 ProcessCPUStats `json:"cpu"`
	Mem                 ProcessMemStats `json:"mem"`
}

type ProcessCPUStats struct {
	Percent       *int64 `json:"percent"`
	TotalInMillis *int64 `json:"total_in_millis"`
}

type ProcessMemStats struct {
	TotalVirtualInBytes *int64 `json:"total_virtual_in_bytes"`
}

type Jvm struct {
	Timestamp       *int64          `json:"timestamp"`
	UptimeInMillis  *int64          `json:"uptime_in_millis"`
	JvmMemStats     JvmMemStats     `json:"mem"`
	JvmThreadsStats JvmThreadsStats `json:"threads"`
	JvmGcStats      JvmGcStats      `json:"gc"`
	BufferPools     BufferPools     `json:"buffer_pools"`
	Classes         Classes         `json:"classes"`
}

type JvmMemStats struct {
	HeapUsedInBytes         *int64 `json:"heap_used_in_bytes"`
	HeapUsedPercent         *int64 `json:"heap_used_percent"`
	HeapCommittedInBytes    *int64 `json:"heap_committed_in_bytes"`
	HeapMaxInBytes          *int64 `json:"heap_max_in_bytes"`
	NonHeapUsedInBytes      *int64 `json:"non_heap_used_in_bytes"`
	NonHeapCommittedInBytes *int64 `json:"non_heap_committed_in_bytes"`
	Pools                   Pools  `json:"pools"`
}

type JvmThreadsStats struct {
	Count     *int64 `json:"count"`
	PeakCount *int64 `json:"peak_count"`
}

type JvmGcStats struct {
	Collectors Collectors `json:"collectors"`
}

type Collectors struct {
	Young Young `json:"young"`
	Old   Old   `json:"old"`
}

type Young struct {
	CollectionCount        *int64 `json:"collection_count"`
	CollectionTimeInMillis *int64 `json:"collection_time_in_millis"`
}

type Old struct {
	CollectionCount        *int64 `json:"collection_count"`
	CollectionTimeInMillis *int64 `json:"collection_time_in_millis"`
}

type Pools struct {
	Young    PoolsYoungStats    `json:"young"`
	Survivor PoolsSurvivorStats `json:"survivor"`
	Old      PoolsOldStats      `json:"old"`
}

type PoolsYoungStats struct {
	UsedInBytes     *int64 `json:"used_in_bytes"`
	MaxInBytes      *int64 `json:"max_in_bytes"`
	PeakUsedInBytes *int64 `json:"peak_used_in_bytes"`
	PeakMaxInBytes  *int64 `json:"peak_max_in_bytes"`
}

type PoolsSurvivorStats struct {
	UsedInBytes     *int64 `json:"used_in_bytes"`
	MaxInBytes      *int64 `json:"max_in_bytes"`
	PeakUsedInBytes *int64 `json:"peak_used_in_bytes"`
	PeakMaxInBytes  *int64 `json:"peak_max_in_bytes"`
}

type PoolsOldStats struct {
	UsedInBytes     *int64 `json:"used_in_bytes"`
	MaxInBytes      *int64 `json:"max_in_bytes"`
	PeakUsedInBytes *int64 `json:"peak_used_in_bytes"`
	PeakMaxInBytes  *int64 `json:"peak_max_in_bytes"`
}

type BufferPools struct {
	Mapped Mapped `json:"mapped"`
	Direct Direct `json:"direct"`
}

type Mapped struct {
	Count                *int64 `json:"count"`
	UsedInBytes          *int64 `json:"used_in_bytes"`
	TotalCapacityInBytes *int64 `json:"total_capacity_in_bytes"`
}

type Direct struct {
	Count                *int64 `json:"count"`
	UsedInBytes          *int64 `json:"used_in_bytes"`
	TotalCapacityInBytes *int64 `json:"total_capacity_in_bytes"`
}

type Classes struct {
	CurrentLoadedCount *int64 `json:"current_loaded_count"`
	TotalLoadedCount   *int64 `json:"total_loaded_count"`
	TotalUnloadedCount *int64 `json:"total_unloaded_count"`
}

type Transport struct {
	ServerOpen    *int64 `json:"server_open"`
	RxCount       *int64 `json:"rx_count"`
	RxSizeInBytes *int64 `json:"rx_size_in_bytes"`
	TxCount       *int64 `json:"tx_count"`
	TxSizeInBytes *int64 `json:"tx_size_in_bytes"`
}

type HTTP struct {
	CurrentOpen *int64 `json:"current_open"`
	TotalOpened *int64 `json:"total_opened"`
}

type Attributes struct {
	MlMachineMemory *string `json:"ml.machine_memory"`
	XpackInstalled  *string `json:"xpack.installed"`
	MlMaxOpenJobs   *string `json:"ml.max_open_jobs"`
	MlEnabled       *string `json:"ml.enabled"`
}
