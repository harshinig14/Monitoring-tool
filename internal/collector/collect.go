package collector

import (
	"os"
	"runtime"
	"time"

	"MONITORING-TOOL/internal/models"
)

func CollectMetrics() (*models.Metrics, error) {

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	cpuPercent, err := GetCPUPercent()
	if err != nil {
		return nil, err
	}

	memoryStats, err := GetMemoryStats()
	if err != nil {
		return nil, err
	}

	diskStats, err := GetDiskStats()
	if err != nil {
		return nil, err
	}

	networkStats, err := GetNetworkStats()
	if err != nil {
		return nil, err
	}

	openPorts, err := GetOpenPortsCount()
	if err != nil {
		return nil, err
	}

	ioStats, err := GetIOStats()
	if err != nil {
		return nil, err
	}

	_, _, _ = GetTopProcesses()

	metrics := &models.Metrics{
		Timestamp: time.Now().Format(time.RFC3339),
		Hostname:  hostname,
		OS:        runtime.GOOS,

		CPUPercent: cpuPercent,

		MemoryTotal:       memoryStats.Total,
		MemoryUsed:        memoryStats.Used,
		MemoryFree:        memoryStats.Free,
		MemoryUsedPercent: memoryStats.UsedPercent,

		DiskTotal:       diskStats.Total,
		DiskUsed:        diskStats.Used,
		DiskFree:        diskStats.Free,
		DiskUsedPercent: diskStats.UsedPercent,

		BytesSent: networkStats.BytesSent,
		BytesRecv: networkStats.BytesRecv,

		PacketsSent: networkStats.PacketsSent,
		PacketsRecv: networkStats.PacketsRecv,

		OpenPorts: openPorts,

		ReadBytes:  ioStats.ReadBytes,
		WriteBytes: ioStats.WriteBytes,
	}

	return metrics, nil
}