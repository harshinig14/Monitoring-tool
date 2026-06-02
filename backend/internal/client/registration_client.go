package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"runtime"

	"MONITORING-TOOL/internal/models"
)

func RegisterDevice(serverURL string) (int, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
	}

	currentUser, err := user.Current()
	username := "unknown-user"
	if err == nil {
		username = currentUser.Username
	}

	osType := runtime.GOOS

	reqData := models.RegisterRequest{
		Username:    username,
		MachineName: hostname,
		OSType:      osType,
	}

	payload, _ := json.Marshal(reqData)

	resp, err := http.Post(fmt.Sprintf("%s/api/v1/register", serverURL), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("server returned status: %s", resp.Status)
	}

	var resData models.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&resData); err != nil {
		return 0, err
	}

	return resData.UserID, nil
}
