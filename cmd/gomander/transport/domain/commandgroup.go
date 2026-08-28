package domain

import (
	commandgroupdomain "gomander/internal/commandgroup/domain"
)

type CommandGroup struct {
	Id        string    `json:"id"`
	ProjectId string    `json:"projectId"`
	Name      string    `json:"name"`
	Commands  []Command `json:"commands"`
	Position  int       `json:"position"`
}

func FromCommandGroup(commandGroup commandgroupdomain.CommandGroup) CommandGroup {
	return CommandGroup{
		Id:        commandGroup.Id,
		ProjectId: commandGroup.ProjectId,
		Name:      commandGroup.Name,
		Commands:  FromCommands(commandGroup.Commands),
		Position:  commandGroup.Position,
	}
}

func FromCommandGroups(commandGroups []commandgroupdomain.CommandGroup) []CommandGroup {
	return mapSlice(commandGroups, FromCommandGroup)
}

func (g CommandGroup) ToDomain() commandgroupdomain.CommandGroup {
	return commandgroupdomain.CommandGroup{
		Id:        g.Id,
		ProjectId: g.ProjectId,
		Name:      g.Name,
		Commands:  mapSlice(g.Commands, Command.ToDomain),
		Position:  g.Position,
	}
}
