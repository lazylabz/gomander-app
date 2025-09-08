package app

import (
	"encoding/json"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type CommandGroupJSONv1 struct {
	Name       string   `json:"name"`
	CommandIds []string `json:"commandIds"`
}

type CommandJSONv1 struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Command          string `json:"command"`
	WorkingDirectory string `json:"workingDirectory"`
}

type ProjectExportJSONv1 struct {
	Version       int                  `json:"version"`
	Name          string               `json:"name"`
	Commands      []CommandJSONv1      `json:"commands"`
	CommandGroups []CommandGroupJSONv1 `json:"commandGroups"`
}

func (a *App) GetProjectToImport() (projectJSON *ProjectExportJSONv1, err error) {
	filePath, err := a.runtimeFacade.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Select a project file",
		Filters: []runtime.FileFilter{{DisplayName: "JSON Files", Pattern: "*.json"}},
	})
	if err != nil {
		return nil, err
	}

	if filePath == "" {
		return nil, nil // User canceled
	}

	// Read entire file into memory
	fileData, err := a.fsFacade.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON data
	err = json.Unmarshal(fileData, &projectJSON)
	if err != nil {
		return nil, err
	}

	return
}
