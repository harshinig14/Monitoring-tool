package collector

import (
	"github.com/shirou/gopsutil/v3/disk"
)

type IOStats struct {
	ReadBytes  uint64
	WriteBytes uint64
}

func GetIOStats() (*IOStats, error) {

	counters, err := disk.IOCounters()
	if err != nil {
		return nil, err
	}

	var totalRead uint64
	var totalWrite uint64

	for _, counter := range counters {
		totalRead += counter.ReadBytes
		totalWrite += counter.WriteBytes
	}

	return &IOStats{
		ReadBytes:  totalRead,
		WriteBytes: totalWrite,
	}, nil
}