package collector

import (
	"sort"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessInfo struct {
	PID        int32
	Name       string
	CPUPercent float64
	MemoryMB   float32
}

func GetTopProcesses() ([]ProcessInfo, []ProcessInfo, error) {

	processes, err := process.Processes()
	if err != nil {
		return nil, nil, err
	}

	var cpuProcesses []ProcessInfo
	var memoryProcesses []ProcessInfo

	for _, p := range processes {

		name, _ := p.Name()

		cpuPercent, _ := p.CPUPercent()

		memInfo, _ := p.MemoryInfo()

		var memoryMB float32

		if memInfo != nil {
			memoryMB = float32(memInfo.RSS) / 1024 / 1024
		}

		info := ProcessInfo{
			PID:        p.Pid,
			Name:       name,
			CPUPercent: cpuPercent,
			MemoryMB:   memoryMB,
		}

		cpuProcesses = append(cpuProcesses, info)
		memoryProcesses = append(memoryProcesses, info)
	}

	sort.Slice(cpuProcesses, func(i, j int) bool {
		return cpuProcesses[i].CPUPercent > cpuProcesses[j].CPUPercent
	})

	sort.Slice(memoryProcesses, func(i, j int) bool {
		return memoryProcesses[i].MemoryMB > memoryProcesses[j].MemoryMB
	})

	if len(cpuProcesses) > 5 {
		cpuProcesses = cpuProcesses[:5]
	}

	if len(memoryProcesses) > 5 {
		memoryProcesses = memoryProcesses[:5]
	}

	return cpuProcesses, memoryProcesses, nil
}