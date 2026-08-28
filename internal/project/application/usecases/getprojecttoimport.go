package usecases

import (
	"encoding/json"
	"path"

	"github.com/google/uuid"

	"gomander/internal/dialog"
	"gomander/internal/facade"
	projectdomain "gomander/internal/project/domain"
)

type FileType string

const (
	FileTypeGomander    FileType = "gomander_export"
	FileTypePackageJSON FileType = "package_json"
)

type GetProjectToImport struct {
	dialogs  dialog.Dialogs
	fsFacade facade.FsFacade
}

func NewGetProjectToImport(
	dialogs dialog.Dialogs,
	fsFacade facade.FsFacade,
) *GetProjectToImport {
	return &GetProjectToImport{
		dialogs:  dialogs,
		fsFacade: fsFacade,
	}
}

func (uc *GetProjectToImport) Execute(fileType FileType) (*projectdomain.ProjectExportJSONv1, error) {
	var projectJSON *projectdomain.ProjectExportJSONv1

	filePath, err := uc.dialogs.AskForFileToOpen(OpenFileRequestByFileType[fileType])
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

	processor := ProcessorsByFileType[fileType]

	projectJSON, err = processor(fileData, filePath)
	if err != nil {
		return nil, err
	}

	return projectJSON, nil
}

var OpenFileRequestByFileType = map[FileType]dialog.OpenFileRequest{
	FileTypeGomander: {
		Title:   "Select an exported Gomander project file",
		Filters: []dialog.FileFilter{{DisplayName: "JSON Files", Pattern: "*.json"}},
	},
	FileTypePackageJSON: {
		Title:   "Select a package.json file",
		Filters: []dialog.FileFilter{{DisplayName: "package.json", Pattern: "*.json"}},
	},
}

var ProcessorsByFileType = map[FileType]func([]byte, string) (*projectdomain.ProjectExportJSONv1, error){
	FileTypeGomander:    parseGomanderExportedProject,
	FileTypePackageJSON: parsePackageJSON,
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

	folderPath := path.Dir(filePath)

	var projectExport = &projectdomain.ProjectExportJSONv1{
		Version:          1,
		Name:             packageJSON.Name,
		Commands:         make([]projectdomain.CommandJSONv1, 0),
		CommandGroups:    make([]projectdomain.CommandGroupJSONv1, 0),
		WorkingDirectory: folderPath,
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
