package project

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const DesiredProjectJSONVersion = 1

const ProjectsFolder = "projects"

type ProjectPersistedJSON struct {
	*Project
	Version int `json:"version"`
}

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

func LoadProject(projectConfigId string) (p *Project, err error) {
	file, err := findOrCreateProjectConfigFile(projectConfigId)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := file.Close()
		if err == nil {
			err = closeErr
		}
	}()

	pj := &ProjectPersistedJSON{}

	// Read the config from the file
	decoder := json.NewDecoder(file)
	err = decoder.Decode(pj)
	if err != nil {
		return nil, err
	}

	p = ProjectFromJSON(pj)

	return
}

func DeleteProject(projectConfigId string) error {
	exists, err := projectFileExists(projectConfigId)

	if err != nil {
		return err
	}

	if !exists {
		return errors.New("project not found")
	}

	filePath, err := getProjectPath(projectConfigId)
	if err != nil {
		return err
	}

	// Remove the file
	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func SaveProject(config *Project) (err error) {
	file, err := findOrCreateProjectConfigFile(config.Id)
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

	pj := JSONFromProject(config)

	// Save the config to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(pj)
	if err != nil {
		return err
	}

	return
}

func ExportProject(project *Project, exportPath string) (err error) {
	// Ensure the export directory exists
	err = os.MkdirAll(filepath.Dir(exportPath), 0755)
	if err != nil {
		return err
	}

	// Create or open the export file
	file, err := os.OpenFile(exportPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer func() {
		closeFileErr := file.Close()
		if err == nil {
			err = closeFileErr
		}
	}()

	// Omit fields that should not be exported
	project.BaseWorkingDirectory = ""

	pj := JSONFromProject(project)

	// Write the project config to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(pj)
	if err != nil {
		return err
	}

	return
}

func ImportProject(project Project) (err error) {
	// Check if there is a project with the same ID. If so, generate a new UUID for the project.
	exists, err := projectFileExists(project.Id)
	if err != nil {
		return err
	}

	if exists {
		project.Id = uuid.New().String()
	}

	// Save the imported project
	err = SaveProject(&project)
	if err != nil {
		return err
	}

	return
}

func LoadProjectFromPath(filePath string) (project *Project, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeFileErr := file.Close()
		if err == nil {
			err = closeFileErr
		}
	}()

	pj := &ProjectPersistedJSON{}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(pj)
	if err != nil {
		return nil, err
	}

	project = ProjectFromJSON(pj)

	return
}

func projectFileExists(projectId string) (bool, error) {
	projectFilePath, err := getProjectPath(projectId)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(projectFilePath)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // File does not exist
		}
		return false, err // Some other error occurred
	}

	return true, nil // File exists
}

func findOrCreateProjectConfigFile(projectConfigId string) (*os.File, error) {
	folderPath, err := getProjectFolderPath()
	if err != nil {
		return nil, err
	}
	filePath, err := getProjectPath(projectConfigId)
	if err != nil {
		return nil, err
	}

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

func getProjectFolderPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	folderPath := filepath.Join(configDir, "gomander", ProjectsFolder)

	return folderPath, nil
}

func getProjectPath(projectId string) (string, error) {
	folderPath, err := getProjectFolderPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(folderPath, projectId+".json"), nil
}

func ProjectFromJSON(pj *ProjectPersistedJSON) (p *Project) {
	if pj.Version != DesiredProjectJSONVersion {
		// TODO: Execute migration logic when needed
		panic(errors.New("project version mismatch"))
	}

	p = &Project{
		Id:                   pj.Id,
		Name:                 pj.Name,
		BaseWorkingDirectory: pj.BaseWorkingDirectory,
		Commands:             pj.Commands,
		CommandGroups:        pj.CommandGroups,
		OrderedCommandIds:    pj.OrderedCommandIds,
	}

	return
}

func JSONFromProject(p *Project) (pj *ProjectPersistedJSON) {
	pj = &ProjectPersistedJSON{
		Project: p,
		Version: DesiredProjectJSONVersion,
	}

	return
}
