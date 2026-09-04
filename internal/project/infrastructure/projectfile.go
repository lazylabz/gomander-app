package infrastructure

import (
	"encoding/json"
	"fmt"

	"gomander/internal/facade"
	"gomander/internal/project/domain"
)

// ProjectFile is the file Gomander exports a Project to and imports one from.
// Write always produces version 2. Read accepts version 1, which cannot carry
// a Command's link or error patterns, and version 2, which can. An older file
// therefore imports as "no link, no patterns" — what those Commands were.
type ProjectFile struct {
	fs facade.FsFacade
	v1 *ProjectFileV1
	v2 *ProjectFileV2
}

func NewProjectFile(fs facade.FsFacade) *ProjectFile {
	return &ProjectFile{
		fs: fs,
		v1: NewProjectFileV1(fs),
		v2: NewProjectFileV2(fs),
	}
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
	case 0, 1:
		return f.v1.Read(filePath)
	case 2:
		return decodeProjectV2(data)
	default:
		return nil, fmt.Errorf("unsupported project file version %d", header.Version)
	}
}

func (f *ProjectFile) Write(filePath string, blueprint domain.Blueprint) error {
	return f.v2.Write(filePath, blueprint)
}
