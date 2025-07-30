package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type UserConfig struct {
	ExtraPaths          []string `json:"extra_paths"`
	LastOpenedProjectId string   `json:"last_opened_project_id"`
}

func EmptyUserConfig() *UserConfig {
	return &UserConfig{
		ExtraPaths:          make([]string, 0),
		LastOpenedProjectId: "",
	}
}

func LoadUserConfig() (*UserConfig, error) {
	file, err := findOrCreateUserConfigFile()
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	stat, err := os.Stat(file.Name())
	if err != nil {
		return nil, err
	}

	if stat.Size() == 0 {
		err := SaveUserConfig(EmptyUserConfig())
		if err != nil {
			return nil, err
		}
		return EmptyUserConfig(), nil
	}

	// Read the config from the file
	var config UserConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveUserConfig(config *UserConfig) error {
	file, err := findOrCreateUserConfigFile()
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// Truncate the file to ensure clean write
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	// Save the config to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}

func findOrCreateUserConfigFile() (*os.File, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	folderPath := filepath.Join(configDir, "gomander")
	filePath := filepath.Join(folderPath, "user_config.json")

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
