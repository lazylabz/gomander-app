package main

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	Commands map[string]Command `json:"commands"`
}

func loadConfig() (*Config, error) {
	file, err := findOrCreateConfigFile()
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err = file.Close()
	}(file)

	// Read the config from the file
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, err
}

func saveConfig(config *Config) error {
	file, err := findOrCreateConfigFile()
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err = file.Close()
	}(file)

	// Save the config to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(config)

	return err
}

func findOrCreateConfigFile() (*os.File, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := strings.Builder{}

	path.WriteString(configDir)
	path.WriteString(string(os.PathSeparator))
	path.WriteString("gomander")

	folderPath := path.String()

	path.WriteString(string(os.PathSeparator))
	path.WriteString("settings.json")

	filePath := path.String()

	// Ensure the directory exists
	err = os.MkdirAll(folderPath, 0755)
	if err != nil {
		return nil, err
	}

	// Open the file, creating it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return nil, err
	}

	return file, nil
}
