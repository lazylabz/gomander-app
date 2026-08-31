package domain

import (
	commanddomain "gomander/internal/command/domain"
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

// FromCommandGroup narrows the Commands the read path still carries to the ids
// the wire speaks in. It becomes a copy once that read path names them too.
func FromCommandGroup(commandGroup commandgroupdomain.CommandGroup) CommandGroup {
	return CommandGroup{
		Id:         commandGroup.Id,
		ProjectId:  commandGroup.ProjectId,
		Name:       commandGroup.Name,
		CommandIds: mapSlice(commandGroup.Commands, func(command commanddomain.Command) string { return command.Id }),
		Position:   commandGroup.Position,
	}
}

func FromCommandGroups(commandGroups []commandgroupdomain.CommandGroup) []CommandGroup {
	return mapSlice(commandGroups, FromCommandGroup)
}

func (g CommandGroup) ToDomain() commandgroupdomain.CommandGroupWithCommandIds {
	return commandgroupdomain.CommandGroupWithCommandIds{
		Id:         g.Id,
		ProjectId:  g.ProjectId,
		Name:       g.Name,
		CommandIds: g.CommandIds,
		Position:   g.Position,
	}
}
