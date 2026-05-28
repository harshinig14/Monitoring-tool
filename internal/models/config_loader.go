package models

import (
	"encoding/json"
	"os"
)

func LoadConfig(path string) (*Config, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var config Config

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}