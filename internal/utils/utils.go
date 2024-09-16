package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the configuration structure
type Config struct {
	Distributors struct {
		DigiKey struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		} `json:"digikey"`
		Mouser struct {
			APIKey string `json:"api_key"`
		} `json:"mouser"`
	} `json:"distributors"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %v", err)
	}

	return config, nil
}
