package domain

import (
	commandgroupdomain "gomander/internal/commandgroup/domain"
)

// CommandGroup names the Commands the Command Group holds, in the order it
// holds them. The frontend already loads every Command of the opened Project,
// so it resolves the ids against those rather than being sent the same Commands
// a second time.
type CommandGroup struct {
	Id         string   `json:"id"`
	ProjectId  string   `json:"projectId"`
	Name       string   `json:"name"`
	CommandIds []string `json:"commandIds"`
	Position   int      `json:"position"`
}

func FromCommandGroup(commandGroup commandgroupdomain.CommandGroup) CommandGroup {
	return CommandGroup{
		Id:         commandGroup.Id,
		ProjectId:  commandGroup.ProjectId,
		Name:       commandGroup.Name,
		CommandIds: commandGroup.CommandIds,
		Position:   commandGroup.Position,
	}
}

func FromCommandGroups(commandGroups []commandgroupdomain.CommandGroup) []CommandGroup {
	return mapSlice(commandGroups, FromCommandGroup)
}

func (g CommandGroup) ToDomain() commandgroupdomain.CommandGroup {
	return commandgroupdomain.CommandGroup{
		Id:         g.Id,
		ProjectId:  g.ProjectId,
		Name:       g.Name,
		CommandIds: g.CommandIds,
		Position:   g.Position,
	}
}
