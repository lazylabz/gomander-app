package project

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ProjectsFolder = "projects"

func GetAllProjectsAvailableInProjectsFolder() ([]*Project, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	folderPath := filepath.Join(configDir, "gomander", ProjectsFolder)

	// Ensure the directory exists
	err = os.MkdirAll(folderPath, 0755)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	projects := make([]*Project, 0, len(files))

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			projectConfigId := file.Name()[:len(file.Name())-5] // Remove .json extension
			project, err := LoadProject(projectConfigId)
			if err != nil {
				return nil, err
			}
			projects = append(projects, project)
		}
	}

	return projects, nil
}

func LoadProject(projectConfigId string) (*Project, error) {
	file, err := findOrCreateProjectConfigFile(projectConfigId)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		closeErr := file.Close()
		if err == nil {
			err = closeErr
		}
	}(file)

	// Read the config from the file
	var config Project
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func DeleteProject(projectConfigId string) error {
	file, err := findOrCreateProjectConfigFile(projectConfigId)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		closeErr := file.Close()
		if err == nil {
			err = closeErr
		}
	}(file)

	// Remove the file
	err = os.Remove(file.Name())
	if err != nil {
		return err
	}

	return nil
}

func SaveProject(config *Project) error {
	file, err := findOrCreateProjectConfigFile(config.Id)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		closeErr := file.Close()
		if err == nil {
			err = closeErr
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

func ExportProject(project *Project, exportPath string) error {
	// Ensure the export directory exists
	err := os.MkdirAll(filepath.Dir(exportPath), 0755)
	if err != nil {
		return err
	}

	// Create or open the export file
	file, err := os.OpenFile(exportPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		closeFileErr := file.Close()
		if err != nil && closeFileErr != nil {
			err = closeFileErr
		}
	}(file)

	// Write the project config to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(project)
	if err != nil {
		return err
	}

	return nil
}

func findOrCreateProjectConfigFile(projectConfigId string) (*os.File, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	folderPath := filepath.Join(configDir, "gomander", ProjectsFolder)
	filePath := filepath.Join(folderPath, projectConfigId+".json")

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
