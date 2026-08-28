package main

import (
	"gomander/internal/helpers/array"
	projectdomain "gomander/internal/project/domain"
)

// ProjectBlueprint is what the frontend sees of a projectdomain.Blueprint. The
// domain type carries no serialization tags, so the names on the wire are
// decided here, at the boundary that owns them.
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

func newProjectBlueprint(blueprint projectdomain.Blueprint) ProjectBlueprint {
	return ProjectBlueprint{
		Name:             blueprint.Name,
		WorkingDirectory: blueprint.WorkingDirectory,
		Commands: array.Map(blueprint.Commands, func(cmd projectdomain.BlueprintCommand) CommandBlueprint {
			return CommandBlueprint{
				Id:               cmd.Id,
				Name:             cmd.Name,
				Command:          cmd.Command,
				WorkingDirectory: cmd.WorkingDirectory,
			}
		}),
		CommandGroups: array.Map(blueprint.CommandGroups, func(group projectdomain.BlueprintCommandGroup) CommandGroupBlueprint {
			return CommandGroupBlueprint{
				Id:         group.Id,
				Name:       group.Name,
				CommandIds: group.CommandIds,
			}
		}),
	}
}

func (b ProjectBlueprint) toDomain() projectdomain.Blueprint {
	return projectdomain.Blueprint{
		Name:             b.Name,
		WorkingDirectory: b.WorkingDirectory,
		Commands: array.Map(b.Commands, func(cmd CommandBlueprint) projectdomain.BlueprintCommand {
			return projectdomain.BlueprintCommand{
				Id:               cmd.Id,
				Name:             cmd.Name,
				Command:          cmd.Command,
				WorkingDirectory: cmd.WorkingDirectory,
			}
		}),
		CommandGroups: array.Map(b.CommandGroups, func(group CommandGroupBlueprint) projectdomain.BlueprintCommandGroup {
			return projectdomain.BlueprintCommandGroup{
				Id:         group.Id,
				Name:       group.Name,
				CommandIds: group.CommandIds,
			}
		}),
	}
}
