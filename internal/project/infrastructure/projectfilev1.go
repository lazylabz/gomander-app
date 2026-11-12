package infrastructure

import (
	"encoding/json"

	"gomander/internal/project/domain"
)

// The shapes below are version 1 of the file Gomander exports a Project to, from
// before a command carried a link and error patterns. It is no longer written;
// keeping it byte for byte is what lets a file exported by an older Gomander
// still be imported. A newer version is a new codec, not an edit of this one.
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

func decodeProjectV1(data []byte) (*domain.Blueprint, error) {
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
