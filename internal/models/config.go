package models

type Config struct {
	CollectionIntervalSeconds int    `json:"collection_interval_seconds"`
	TraceDirectory            string `json:"trace_directory"`

	EnableCPU               bool `json:"enable_cpu"`
	EnableMemory            bool `json:"enable_memory"`
	EnableDisk              bool `json:"enable_disk"`
	EnableNetwork           bool `json:"enable_network"`
	EnablePorts             bool `json:"enable_ports"`
	EnableIO                bool `json:"enable_io"`
	EnableProcessMonitoring bool `json:"enable_process_monitoring"`
}