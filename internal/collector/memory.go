package collector

import (
	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryStats struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

func GetMemoryStats() (*MemoryStats, error) {

	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &MemoryStats{
		Total:       vm.Total,
		Used:        vm.Used,
		Free:        vm.Free,
		UsedPercent: vm.UsedPercent,
	}, nil
}