package infrastructure

import (
	"encoding/json"

	"gomander/internal/facade"
	"gomander/internal/project/domain"
)

// ProjectFileV2 is version 2 of the file Gomander exports a Project to. It
// carries each Command's link and error patterns, which version 1 had nowhere
// to put. It is the version written today; version 1 remains the reader for
// files an older Gomander already produced.
type ProjectFileV2 struct {
	fs facade.FsFacade
}

func NewProjectFileV2(fs facade.FsFacade) *ProjectFileV2 {
	return &ProjectFileV2{fs: fs}
}

type projectV2 struct {
	Version          int              `json:"version"`
	Name             string           `json:"name"`
	WorkingDirectory string           `json:"workingDirectory"`
	Commands         []commandV2      `json:"commands"`
	CommandGroups    []commandGroupV2 `json:"commandGroups"`
}

type commandV2 struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	Command          string   `json:"command"`
	WorkingDirectory string   `json:"workingDirectory"`
	Link             string   `json:"link"`
	ErrorPatterns    []string `json:"errorPatterns"`
}

type commandGroupV2 struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	CommandIds []string `json:"commandIds"`
}

func (f *ProjectFileV2) Read(filePath string) (*domain.Blueprint, error) {
	data, err := f.fs.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return decodeProjectV2(data)
}

func decodeProjectV2(data []byte) (*domain.Blueprint, error) {
	var file projectV2
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
			Link:             command.Link,
			ErrorPatterns:    emptyStrings(command.ErrorPatterns),
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

func (f *ProjectFileV2) Write(filePath string, blueprint domain.Blueprint) error {
	file := projectV2{
		Version:          2,
		Name:             blueprint.Name,
		WorkingDirectory: blueprint.WorkingDirectory,
	}

	for _, command := range blueprint.Commands {
		file.Commands = append(file.Commands, commandV2{
			Id:               command.Id,
			Name:             command.Name,
			Command:          command.Command,
			WorkingDirectory: command.WorkingDirectory,
			Link:             command.Link,
			ErrorPatterns:    emptyStrings(command.ErrorPatterns),
		})
	}

	for _, group := range blueprint.CommandGroups {
		file.CommandGroups = append(file.CommandGroups, commandGroupV2{
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

// emptyStrings turns a missing list into an empty one so a file never carries
// null for error patterns, and a Blueprint never has to distinguish the two.
func emptyStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}
