package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = "/.gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find .gatorconfig.json: %w", err)
	}
	filename := homeDir + configFileName
	return filename, nil
}

func Read() (Config, error) {
	var config Config
	filename, err := getConfigFilePath()
	if err != nil {
		return config, err
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshaling json: %w", err)
	}
	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}
	filename, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, jsonData, 0666)
	if err != nil {
		return err
	}
	return nil
}
