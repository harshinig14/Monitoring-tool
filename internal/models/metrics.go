package models

import "strconv"

// Metrics stores complete monitoring information
// collected from the system.
type Metrics struct {

	// Basic Metadata
	Timestamp string
	Hostname  string
	OS        string

	// CPU Metrics
	CPUPercent float64

	// Memory Metrics
	MemoryTotal       uint64
	MemoryUsed        uint64
	MemoryFree        uint64
	MemoryUsedPercent float64

	// Disk Metrics
	DiskTotal       uint64
	DiskUsed        uint64
	DiskFree        uint64
	DiskUsedPercent float64

	// Network Metrics
	BytesSent uint64
	BytesRecv uint64

	PacketsSent uint64
	PacketsRecv uint64

	// Ports
	OpenPorts int

	// I/O Metrics
	ReadBytes  uint64
	WriteBytes uint64
}

// CSVHeaders defines CSV column names.
var CSVHeaders = []string{
	"timestamp",
	"hostname",
	"os",

	"cpu_percent",

	"memory_total",
	"memory_used",
	"memory_free",
	"memory_used_percent",

	"disk_total",
	"disk_used",
	"disk_free",
	"disk_used_percent",

	"bytes_sent",
	"bytes_recv",

	"packets_sent",
	"packets_recv",

	"open_ports",

	"read_bytes",
	"write_bytes",
}

// ToCSVRow converts Metrics into CSV row format.
func (m Metrics) ToCSVRow() []string {
	return []string{
		m.Timestamp,
		m.Hostname,
		m.OS,

		formatFloat(m.CPUPercent),

		formatUint(m.MemoryTotal),
		formatUint(m.MemoryUsed),
		formatUint(m.MemoryFree),
		formatFloat(m.MemoryUsedPercent),

		formatUint(m.DiskTotal),
		formatUint(m.DiskUsed),
		formatUint(m.DiskFree),
		formatFloat(m.DiskUsedPercent),

		formatUint(m.BytesSent),
		formatUint(m.BytesRecv),

		formatUint(m.PacketsSent),
		formatUint(m.PacketsRecv),

		formatInt(m.OpenPorts),

		formatUint(m.ReadBytes),
		formatUint(m.WriteBytes),
	}
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func formatUint(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func formatInt(value int) string {
	return strconv.Itoa(value)
}