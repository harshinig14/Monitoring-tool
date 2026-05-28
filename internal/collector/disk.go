package collector

import (
	"runtime"

	"github.com/shirou/gopsutil/v3/disk"
)

type DiskStats struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

func GetDiskStats() (*DiskStats, error) {

	path := "/"

	if runtime.GOOS == "windows" {
		path = "C:\\"
	}

	usage, err := disk.Usage(path)
	if err != nil {
		return nil, err
	}

	return &DiskStats{
		Total:       usage.Total,
		Used:        usage.Used,
		Free:        usage.Free,
		UsedPercent: usage.UsedPercent,
	}, nil
}