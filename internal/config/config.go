package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DatabaseUrl     string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(homeDir, configFileName)

	return configPath, nil
}

func getConfigFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func write(newConfigData *Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	updatedData, err := json.MarshalIndent(newConfigData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Read() Config {
	var config Config

	configPath, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Error calling home directory: ", err)
		return config
	}

	file, err := getConfigFile(configPath)
	if err != nil {
		fmt.Println("Error: The config file could not be opened: ", err)
		return config
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		fmt.Println("Error: Config JSON data could not be decoded: ", err)
		return config
	}

	return config
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	err := write(c)
	if err != nil {
		return err
	}
	return nil
}
