package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gomander/internal/facade"
	projectdomain "gomander/internal/project/domain"
)

type FileType string

const (
	FileTypeGomander FileType = "gomander_export"
	PackageJSON      FileType = "package_json"
)

type GetProjectToImport interface {
	Execute(fileType FileType) (*projectdomain.ProjectExportJSONv1, error)
}

type DefaultGetProjectToImport struct {
	ctx           context.Context
	runtimeFacade facade.RuntimeFacade
	fsFacade      facade.FsFacade
}

func NewGetProjectToImport(
	ctx context.Context,
	runtimeFacade facade.RuntimeFacade,
	fsFacade facade.FsFacade,
) *DefaultGetProjectToImport {
	return &DefaultGetProjectToImport{
		ctx:           ctx,
		runtimeFacade: runtimeFacade,
		fsFacade:      fsFacade,
	}
}

func (uc *DefaultGetProjectToImport) Execute(fileType FileType) (*projectdomain.ProjectExportJSONv1, error) {
	var projectJSON *projectdomain.ProjectExportJSONv1

	filePath, err := uc.runtimeFacade.OpenFileDialog(uc.ctx, OpenDialogOptionsByFileType[fileType])
	if err != nil {
		return nil, err
	}

	if filePath == "" {
		return nil, nil // User canceled
	}

	// Read entire file into memory
	fileData, err := uc.fsFacade.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	processor, exists := ProcessorsByFileType[fileType]
	if !exists {
		return nil, errors.New(fmt.Sprintf("file type %s is not supported", fileType))
	}

	projectJSON, err = processor(fileData, filePath)
	if err != nil {
		return nil, err
	}

	return projectJSON, nil
}

var OpenDialogOptionsByFileType = map[FileType]runtime.OpenDialogOptions{
	FileTypeGomander: {
		Title:   "Select a exported Gomander project file",
		Filters: []runtime.FileFilter{{DisplayName: "JSON Files", Pattern: "*.json"}},
	},
	PackageJSON: {
		Title:   "Select a package.json file",
		Filters: []runtime.FileFilter{{DisplayName: "package.json", Pattern: "package.json"}},
	},
}

var ProcessorsByFileType = map[FileType]func([]byte, string) (*projectdomain.ProjectExportJSONv1, error){
	FileTypeGomander: parseGomanderExportedProject,
	PackageJSON:      parsePackageJSON,
}

func parseGomanderExportedProject(data []byte, _ string) (*projectdomain.ProjectExportJSONv1, error) {
	var projectJSON *projectdomain.ProjectExportJSONv1
	err := json.Unmarshal(data, &projectJSON)
	if err != nil {
		return nil, err
	}
	return projectJSON, nil
}

func parsePackageJSON(data []byte, filePath string) (*projectdomain.ProjectExportJSONv1, error) {
	var packageJSON struct {
		Name    string            `json:"name"`
		Scripts map[string]string `json:"scripts"`
	}
	err := json.Unmarshal(data, &packageJSON)
	if err != nil {
		return nil, err
	}

	var projectExport = &projectdomain.ProjectExportJSONv1{
		Version:          1,
		Name:             packageJSON.Name,
		Commands:         make([]projectdomain.CommandJSONv1, 0),
		CommandGroups:    make([]projectdomain.CommandGroupJSONv1, 0),
		WorkingDirectory: filePath,
	}

	for scriptName, scriptCmd := range packageJSON.Scripts {
		command := projectdomain.CommandJSONv1{
			Id:               uuid.NewString(),
			Name:             scriptName,
			Command:          scriptCmd,
			WorkingDirectory: "",
		}
		projectExport.Commands = append(projectExport.Commands, command)
	}

	return projectExport, nil
}
