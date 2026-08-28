package domain

import (
	commanddomain "gomander/internal/command/domain"
)

type Command struct {
	Id               string   `json:"id"`
	ProjectId        string   `json:"projectId"`
	Name             string   `json:"name"`
	Command          string   `json:"command"`
	WorkingDirectory string   `json:"workingDirectory"`
	Position         int      `json:"position"`
	Link             string   `json:"link"`
	ErrorPatterns    []string `json:"errorPatterns"`
}

func FromCommand(command commanddomain.Command) Command {
	return Command{
		Id:               command.Id,
		ProjectId:        command.ProjectId,
		Name:             command.Name,
		Command:          command.Command,
		WorkingDirectory: command.WorkingDirectory,
		Position:         command.Position,
		Link:             command.Link,
		ErrorPatterns:    command.ErrorPatterns,
	}
}

func FromCommands(commands []commanddomain.Command) []Command {
	return mapSlice(commands, FromCommand)
}

func (c Command) ToDomain() commanddomain.Command {
	return commanddomain.Command{
		Id:               c.Id,
		ProjectId:        c.ProjectId,
		Name:             c.Name,
		Command:          c.Command,
		WorkingDirectory: c.WorkingDirectory,
		Position:         c.Position,
		Link:             c.Link,
		ErrorPatterns:    c.ErrorPatterns,
	}
}
