package dto

type RealtimeMetrics struct {
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsage  float64 `json:"memory_usage"`
	DiskUsage    float64 `json:"disk_usage"`
	NetworkUsage float64 `json:"network_usage"`
	TraceDate    string  `json:"trace_date"`
}
