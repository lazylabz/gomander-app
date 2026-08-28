package infrastructure

import (
	"encoding/json"
	"path/filepath"

	"gomander/internal/facade"
	"gomander/internal/project/domain"
)

// PackageJSONFile reads an npm manifest as a Project: its scripts become the
// Commands, run from the folder the manifest sits in. Everything this knows
// about npm stays here.
type PackageJSONFile struct {
	fs facade.FsFacade
}

func NewPackageJSONFile(fs facade.FsFacade) *PackageJSONFile {
	return &PackageJSONFile{fs: fs}
}

func (f *PackageJSONFile) Read(filePath string) (*domain.Blueprint, error) {
	data, err := f.fs.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var manifest struct {
		Name    string            `json:"name"`
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	blueprint := domain.Blueprint{
		Name:             manifest.Name,
		WorkingDirectory: filepath.Dir(filePath),
		Commands:         make([]domain.BlueprintCommand, 0, len(manifest.Scripts)),
		CommandGroups:    make([]domain.BlueprintCommandGroup, 0),
	}

	for scriptName, script := range manifest.Scripts {
		blueprint.Commands = append(blueprint.Commands, domain.BlueprintCommand{
			// A script name is unique within the manifest, which is all a
			// Blueprint Id has to be.
			Id:      scriptName,
			Name:    scriptName,
			Command: script,
		})
	}

	return &blueprint, nil
}
