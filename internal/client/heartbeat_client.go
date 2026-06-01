package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"MONITORING-TOOL/internal/models"
)

func SendHeartbeat(userID int, serverURL string) error {
	req := models.HeartbeatRequest{
		UserID: userID,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal heartbeat request: %v", err)
	}

	url := fmt.Sprintf("%s/api/v1/heartbeat", serverURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("http post failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	var res models.HeartbeatResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !res.Success {
		return fmt.Errorf("server returned success: false")
	}

	log.Println("Heartbeat Sent")
	return nil
}
