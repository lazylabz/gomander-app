package infrastructure

import (
	"encoding/json"
	"fmt"

	"gomander/internal/facade"
	"gomander/internal/project/domain"
)

// ProjectFile reads and writes the file Gomander exports a Project to. The file
// names its own version, and each version is a codec of its own so that a file
// an older Gomander exported keeps importing exactly as it did: Read picks the
// codec the file asks for, Write always uses the newest.
type ProjectFile struct {
	fs facade.FsFacade
}

func NewProjectFile(fs facade.FsFacade) *ProjectFile {
	return &ProjectFile{fs: fs}
}

func (f *ProjectFile) Read(filePath string) (*domain.Blueprint, error) {
	data, err := f.fs.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var header struct {
		Version int `json:"version"`
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, err
	}

	switch header.Version {
	case 1:
		return decodeProjectV1(data)
	case 2:
		return decodeProjectV2(data)
	default:
		return nil, fmt.Errorf("unsupported project file version %d", header.Version)
	}
}

func (f *ProjectFile) Write(filePath string, blueprint domain.Blueprint) error {
	data, err := encodeProjectV2(blueprint)
	if err != nil {
		return err
	}

	return f.fs.WriteFile(filePath, data, 0644)
}
