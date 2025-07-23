package config

import (
	"encoding/json"
	"gomander/internal/command"
	"os"
	"strings"
)

type CommandGroup struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	CommandIds []string `json:"commands"`
}

type Config struct {
	Commands      map[string]command.Command `json:"commands"`
	ExtraPaths    []string                   `json:"extra_paths"`
	CommandGroups []CommandGroup             `json:"command_groups"`
}

type UserConfig struct {
	ExtraPaths []string `json:"extraPaths"`
}

func LoadConfigOrPanic() *Config {
	file, err := findOrCreateConfigFile()
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err = file.Close()
	}(file)

	stat, err := os.Stat(file.Name())
	if err != nil {
		panic(err)
	}

	if stat.Size() == 0 {
		return &Config{
			Commands:      make(map[string]command.Command),
			ExtraPaths:    make([]string, 0),
			CommandGroups: make([]CommandGroup, 0),
		}
	}

	// Read the config from the file
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return &config
}

func SaveConfigOrPanic(config *Config) {
	file, err := findOrCreateConfigFile()
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err = file.Close()
	}(file)

	// Truncate the file to ensure clean write
	err = file.Truncate(0)
	if err != nil {
		panic(err)
	}

	// Save the config to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(config)
	if err != nil {
		panic(err)
	}
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
