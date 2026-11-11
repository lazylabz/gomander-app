package infrastructure

import (
	"encoding/json"

	"gomander/internal/project/domain"
)

// Version 2 of the file Gomander exports a Project to, the one written today.
// It differs from version 1 only in what a command carries: its link and error
// patterns.

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
			ErrorPatterns:    command.ErrorPatterns,
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

func encodeProjectV2(blueprint domain.Blueprint) ([]byte, error) {
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
			ErrorPatterns:    command.ErrorPatterns,
		})
	}

	for _, group := range blueprint.CommandGroups {
		file.CommandGroups = append(file.CommandGroups, commandGroupV2{
			Id:         group.Id,
			Name:       group.Name,
			CommandIds: group.CommandIds,
		})
	}

	return json.MarshalIndent(file, "", "  ")
}
