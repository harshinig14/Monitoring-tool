package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportLog struct {
	Period  string `json:"period"`
	AvgCPU  string `json:"avgCpu"`
	AvgMem  string `json:"avgMem"`
	AvgDisk string `json:"avgDisk"`
	AvgNet  string `json:"avgNet"`
}

type ReportsResponse struct {
	AvgCPU       int         `json:"avgCpu"`
	AvgMem       int         `json:"avgMem"`
	AvgDisk      int         `json:"avgDisk"`
	AvgNet       string      `json:"avgNet"`
	ChartLabels  []string    `json:"chartLabels"`
	CPUData      []int       `json:"cpuData"`
	MemoryData   []int       `json:"memoryData"`
	DiskData     []int       `json:"diskData"`
	NetworkData  []int       `json:"networkData"`
	Logs         []ReportLog `json:"logs"`
}

func GetDailyReports(c *gin.Context) {
	resp := ReportsResponse{
		AvgCPU:      46,
		AvgMem:      62,
		AvgDisk:     75,
		AvgNet:      "5.6 MB/s",
		ChartLabels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
		CPUData:     []int{40, 48, 47, 48, 52, 38, 45},
		MemoryData:  []int{60, 62, 61, 64, 68, 58, 61},
		DiskData:    []int{73, 74, 74, 75, 75, 74, 75},
		NetworkData: []int{25, 33, 49, 61, 65, 42, 58},
		Logs: []ReportLog{
			{Period: "Sunday, May 31", AvgCPU: "45%", AvgMem: "61%", AvgDisk: "75%", AvgNet: "5.8 MB/s"},
			{Period: "Saturday, May 30", AvgCPU: "38%", AvgMem: "58%", AvgDisk: "74%", AvgNet: "4.2 MB/s"},
			{Period: "Friday, May 29", AvgCPU: "52%", AvgMem: "68%", AvgDisk: "75%", AvgNet: "7.1 MB/s"},
			{Period: "Thursday, May 28", AvgCPU: "48%", AvgMem: "64%", AvgDisk: "75%", AvgNet: "6.0 MB/s"},
			{Period: "Wednesday, May 27", AvgCPU: "47%", AvgMem: "59%", AvgDisk: "73%", AvgNet: "4.9 MB/s"},
		},
	}
	c.JSON(http.StatusOK, resp)
}

func GetWeeklyReports(c *gin.Context) {
	resp := ReportsResponse{
		AvgCPU:      44,
		AvgMem:      60,
		AvgDisk:     74,
		AvgNet:      "5.2 MB/s",
		ChartLabels: []string{"Week 18", "Week 19", "Week 20", "Week 21", "Week 22"},
		CPUData:     []int{43, 45, 42, 44, 46},
		MemoryData:  []int{59, 61, 59, 60, 62},
		DiskData:    []int{73, 73, 74, 74, 75},
		NetworkData: []int{20, 25, 23, 27, 33},
		Logs: []ReportLog{
			{Period: "Week 22 (May 25 - May 31)", AvgCPU: "46%", AvgMem: "62%", AvgDisk: "75%", AvgNet: "5.6 MB/s"},
			{Period: "Week 21 (May 18 - May 24)", AvgCPU: "44%", AvgMem: "60%", AvgDisk: "74%", AvgNet: "5.1 MB/s"},
			{Period: "Week 20 (May 11 - May 17)", AvgCPU: "42%", AvgMem: "59%", AvgDisk: "74%", AvgNet: "4.8 MB/s"},
		},
	}
	c.JSON(http.StatusOK, resp)
}

func GetMonthlyReports(c *gin.Context) {
	resp := ReportsResponse{
		AvgCPU:      43,
		AvgMem:      59,
		AvgDisk:     74,
		AvgNet:      "5.0 MB/s",
		ChartLabels: []string{"Jan", "Feb", "Mar", "Apr", "May"},
		CPUData:     []int{40, 42, 41, 41, 44},
		MemoryData:  []int{57, 58, 58, 58, 60},
		DiskData:    []int{72, 73, 73, 73, 74},
		NetworkData: []int{18, 20, 22, 23, 25},
		Logs: []ReportLog{
			{Period: "May 2026", AvgCPU: "44%", AvgMem: "60%", AvgDisk: "74%", AvgNet: "5.2 MB/s"},
			{Period: "April 2026", AvgCPU: "41%", AvgMem: "58%", AvgDisk: "73%", AvgNet: "4.9 MB/s"},
		},
	}
	c.JSON(http.StatusOK, resp)
}
