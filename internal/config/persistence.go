package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const DesiredUserConfigJSONVersion = 1

type UserConfigPersistedJSON struct {
	*UserConfig
	Version int `json:"version"`
}

func LoadUserConfig() (c *UserConfig, err error) {
	file, err := findOrCreateUserConfigFile()
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := file.Close()
		if err == nil {
			err = closeErr
		}
	}()

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
	pj := &UserConfigPersistedJSON{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(pj)
	if err != nil {
		return nil, err
	}

	c = UserConfigFromJSON(pj)

	return
}

func SaveUserConfig(config *UserConfig) (err error) {
	file, err := findOrCreateUserConfigFile()
	if err != nil {
		return err
	}

	defer func() {
		closeErr := file.Close()
		if err == nil {
			err = closeErr
		}
	}()

	// Truncate the file to ensure clean write
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	pj := JSONFromUserConfig(config)

	// Save the config to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(pj)
	if err != nil {
		return err
	}

	return
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

func UserConfigFromJSON(pj *UserConfigPersistedJSON) *UserConfig {
	if pj.Version != DesiredUserConfigJSONVersion {
		// TODO: Execute migration logic when needed
		panic(errors.New("user config version mismatch"))
	}

	return &UserConfig{
		EnvironmentPaths:    pj.EnvironmentPaths,
		LastOpenedProjectId: pj.LastOpenedProjectId,
	}
}

func JSONFromUserConfig(config *UserConfig) *UserConfigPersistedJSON {
	return &UserConfigPersistedJSON{
		UserConfig: config,
		Version:    DesiredUserConfigJSONVersion,
	}
}
