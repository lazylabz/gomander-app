package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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

	projectJSON, err = processor(fileData)
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

var ProcessorsByFileType = map[FileType]func([]byte) (*projectdomain.ProjectExportJSONv1, error){
	FileTypeGomander: parseGomanderExportedProject,
	PackageJSON:      parsePackageJSON,
}

func parseGomanderExportedProject(data []byte) (*projectdomain.ProjectExportJSONv1, error) {
	var projectJSON *projectdomain.ProjectExportJSONv1
	err := json.Unmarshal(data, &projectJSON)
	if err != nil {
		return nil, err
	}
	return projectJSON, nil
}

func parsePackageJSON(data []byte) (*projectdomain.ProjectExportJSONv1, error) {
	// Placeholder for actual package.json processing logic
	return nil, nil
}
