package usecases

import (
	"context"
	"encoding/json"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gomander/internal/facade"
	projectdomain "gomander/internal/project/domain"
)

type GetProjectToImport interface {
	Execute() (*projectdomain.ProjectExportJSONv1, error)
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

func (uc *DefaultGetProjectToImport) Execute() (*projectdomain.ProjectExportJSONv1, error) {
	var projectJSON *projectdomain.ProjectExportJSONv1

	filePath, err := uc.runtimeFacade.OpenFileDialog(uc.ctx, runtime.OpenDialogOptions{
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
	fileData, err := uc.fsFacade.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON data
	err = json.Unmarshal(fileData, &projectJSON)
	if err != nil {
		return nil, err
	}

	return projectJSON, nil
}
