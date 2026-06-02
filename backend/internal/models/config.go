package models

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerURL                 string `json:"server_url"`
	CollectionIntervalSeconds int    `json:"collection_interval_seconds"`
	UserID                    int    `json:"user_id"`
	EnableCPU                 bool   `json:"enable_cpu"`
	EnableMemory              bool   `json:"enable_memory"`
	EnableDisk                bool   `json:"enable_disk"`
	EnableNetwork             bool   `json:"enable_network"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &config)
	return &config, err
}

func SaveConfig(path string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}