package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"MONITORING-TOOL/internal/models"
)

func SendMetrics(metrics *models.Metrics, config *models.Config) error {
	networkUsage := float64(metrics.BytesSent + metrics.BytesRecv)

	req := models.MetricsRequest{
		UserID:       config.UserID,
		CPUUsage:     metrics.CPUPercent,
		MemoryUsage:  metrics.MemoryUsedPercent,
		DiskUsage:    metrics.DiskUsedPercent,
		NetworkUsage: networkUsage,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics request: %v", err)
	}

	url := fmt.Sprintf("%s/api/v1/metrics", config.ServerURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("http post failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	var res models.MetricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !res.Success {
		return fmt.Errorf("server returned success: false")
	}

	log.Println("Metrics Uploaded")
	return nil
}
