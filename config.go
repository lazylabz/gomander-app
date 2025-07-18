package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Commands map[string]Command `json:"commands"`
}

func loadConfig() (*Config, error) {
	settingsFile, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(*settingsFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Create a new config file if it does not exist
			defaultConfig := &Config{
				Commands: make(map[string]Command),
			}
			err = saveConfig(defaultConfig)
			if err != nil {
				return nil, err
			}
			file, err = os.Open(*settingsFile)
		}
		return nil, err
	}

	defer func(file *os.File) {
		err = file.Close()
	}(file)

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, err
}

func saveConfig(config *Config) error {
	err := findOrCreateConfigDir()
	if err != nil {
		return err
	}

	settingsFile, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(*settingsFile)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err = file.Close()
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(config)

	return err
}

func findOrCreateConfigDir() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dirName := configDir + string(os.PathSeparator) + "gomander"
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func getConfigFilePath() (*string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	dirName := configDir + string(os.PathSeparator) + "gomander" + string(os.PathSeparator) + "settings.json"

	return &dirName, nil
}
