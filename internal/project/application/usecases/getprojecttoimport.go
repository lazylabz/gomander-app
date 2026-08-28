package usecases

import (
	"fmt"

	"gomander/internal/dialog"
	projectdomain "gomander/internal/project/domain"
)

type FileType string

const (
	FileTypeGomander    FileType = "gomander_export"
	FileTypePackageJSON FileType = "package_json"
)

// BlueprintReader turns a file the user picked into the Project it describes.
// Every format the app can import is one implementation and nothing else.
type BlueprintReader interface {
	Read(filePath string) (*projectdomain.Blueprint, error)
}

type GetProjectToImport struct {
	dialogs dialog.Dialogs
	readers map[FileType]BlueprintReader
}

func NewGetProjectToImport(
	dialogs dialog.Dialogs,
	readers map[FileType]BlueprintReader,
) *GetProjectToImport {
	return &GetProjectToImport{
		dialogs: dialogs,
		readers: readers,
	}
}

func (uc *GetProjectToImport) Execute(fileType FileType) (*projectdomain.Blueprint, error) {
	reader, known := uc.readers[fileType]
	if !known {
		return nil, fmt.Errorf("no reader for file type %q", fileType)
	}

	filePath, err := uc.dialogs.AskForFileToOpen(openFileRequestByFileType[fileType])
	if err != nil {
		return nil, err
	}

	if filePath == "" {
		return nil, nil // User canceled
	}

	return reader.Read(filePath)
}

var openFileRequestByFileType = map[FileType]dialog.OpenFileRequest{
	FileTypeGomander: {
		Title:   "Select an exported Gomander project file",
		Filters: []dialog.FileFilter{{DisplayName: "JSON Files", Pattern: "*.json"}},
	},
	FileTypePackageJSON: {
		Title:   "Select a package.json file",
		Filters: []dialog.FileFilter{{DisplayName: "package.json", Pattern: "*.json"}},
	},
}
