package domain

import (
	"gomander/internal/helpers/array"
	projectdomain "gomander/internal/project/domain"
)

type ProjectBlueprint struct {
	Name             string                  `json:"name"`
	WorkingDirectory string                  `json:"workingDirectory"`
	Commands         []CommandBlueprint      `json:"commands"`
	CommandGroups    []CommandGroupBlueprint `json:"commandGroups"`
}

type CommandBlueprint struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Command          string `json:"command"`
	WorkingDirectory string `json:"workingDirectory"`
}

type CommandGroupBlueprint struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	CommandIds []string `json:"commandIds"`
}

// A Blueprint read from an export file with no commands carries nil, and the
// frontend iterates both lists unguarded, so the mappings below reach for
// array.Map rather than the package's nil-preserving mapSlice: the lists have
// to arrive as [] and never null. Both directions do it, so a Blueprint that
// makes the round trip comes back the shape it left as.
func FromBlueprint(blueprint projectdomain.Blueprint) ProjectBlueprint {
	return ProjectBlueprint{
		Name:             blueprint.Name,
		WorkingDirectory: blueprint.WorkingDirectory,
		Commands:         array.Map(blueprint.Commands, fromBlueprintCommand),
		CommandGroups:    array.Map(blueprint.CommandGroups, fromBlueprintCommandGroup),
	}
}

func (b ProjectBlueprint) ToDomain() projectdomain.Blueprint {
	return projectdomain.Blueprint{
		Name:             b.Name,
		WorkingDirectory: b.WorkingDirectory,
		Commands:         array.Map(b.Commands, CommandBlueprint.ToDomain),
		CommandGroups:    array.Map(b.CommandGroups, CommandGroupBlueprint.ToDomain),
	}
}

func fromBlueprintCommand(command projectdomain.BlueprintCommand) CommandBlueprint {
	return CommandBlueprint{
		Id:               command.Id,
		Name:             command.Name,
		Command:          command.Command,
		WorkingDirectory: command.WorkingDirectory,
	}
}

func (c CommandBlueprint) ToDomain() projectdomain.BlueprintCommand {
	return projectdomain.BlueprintCommand{
		Id:               c.Id,
		Name:             c.Name,
		Command:          c.Command,
		WorkingDirectory: c.WorkingDirectory,
	}
}

func fromBlueprintCommandGroup(group projectdomain.BlueprintCommandGroup) CommandGroupBlueprint {
	return CommandGroupBlueprint{
		Id:         group.Id,
		Name:       group.Name,
		CommandIds: group.CommandIds,
	}
}

func (g CommandGroupBlueprint) ToDomain() projectdomain.BlueprintCommandGroup {
	return projectdomain.BlueprintCommandGroup{
		Id:         g.Id,
		Name:       g.Name,
		CommandIds: g.CommandIds,
	}
}
