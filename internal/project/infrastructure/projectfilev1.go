package infrastructure

import (
	"encoding/json"

	"gomander/internal/facade"
	"gomander/internal/project/domain"
)

// ProjectFileV1 reads and writes the file Gomander exports a Project to. The
// shape below is version 1 of that format and is the only one written; keeping
// it byte for byte is what lets a file exported by an older Gomander still be
// imported. A second version is a second codec, not an edit of this one.
type ProjectFileV1 struct {
	fs facade.FsFacade
}

func NewProjectFileV1(fs facade.FsFacade) *ProjectFileV1 {
	return &ProjectFileV1{fs: fs}
}

type projectV1 struct {
	Version          int              `json:"version"`
	Name             string           `json:"name"`
	WorkingDirectory string           `json:"workingDirectory"`
	Commands         []commandV1      `json:"commands"`
	CommandGroups    []commandGroupV1 `json:"commandGroups"`
}

type commandV1 struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Command          string `json:"command"`
	WorkingDirectory string `json:"workingDirectory"`
}

type commandGroupV1 struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	CommandIds []string `json:"commandIds"`
}

func (f *ProjectFileV1) Read(filePath string) (*domain.Blueprint, error) {
	data, err := f.fs.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var file projectV1
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}

	blueprint := domain.Blueprint{
		Name:             file.Name,
		WorkingDirectory: file.WorkingDirectory,
	}

	for _, command := range file.Commands {
		blueprint.Commands = append(blueprint.Commands, domain.BlueprintCommand{
			Id:               command.Id,
			Name:             command.Name,
			Command:          command.Command,
			WorkingDirectory: command.WorkingDirectory,
		})
	}

	for _, group := range file.CommandGroups {
		blueprint.CommandGroups = append(blueprint.CommandGroups, domain.BlueprintCommandGroup{
			Id:         group.Id,
			Name:       group.Name,
			CommandIds: group.CommandIds,
		})
	}

	return &blueprint, nil
}

func (f *ProjectFileV1) Write(filePath string, blueprint domain.Blueprint) error {
	file := projectV1{
		Version:          1,
		Name:             blueprint.Name,
		WorkingDirectory: blueprint.WorkingDirectory,
	}

	for _, command := range blueprint.Commands {
		file.Commands = append(file.Commands, commandV1{
			Id:               command.Id,
			Name:             command.Name,
			Command:          command.Command,
			WorkingDirectory: command.WorkingDirectory,
		})
	}

	for _, group := range blueprint.CommandGroups {
		file.CommandGroups = append(file.CommandGroups, commandGroupV1{
			Id:         group.Id,
			Name:       group.Name,
			CommandIds: group.CommandIds,
		})
	}

	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}

	return f.fs.WriteFile(filePath, data, 0644)
}
