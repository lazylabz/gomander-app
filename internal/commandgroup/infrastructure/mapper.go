package infrastructure

import (
	commanddomain "gomander/internal/command/domain"
	"gomander/internal/command/infrastructure"
	"gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
)

func ToCommandGroupModel(domainCommandGroup *domain.CommandGroup) CommandGroupModel {
	return CommandGroupModel{
		Id:        domainCommandGroup.Id,
		Name:      domainCommandGroup.Name,
		ProjectId: domainCommandGroup.ProjectId,
		Position:  domainCommandGroup.Position,
	}
}

// ToCommandToCommandGroupModels turns the Command Group's own answer about
// where its Commands sit into the rows that hold it.
func ToCommandToCommandGroupModels(domainCommandGroup *domain.CommandGroup) []CommandToCommandGroupModel {
	return array.Map(
		domainCommandGroup.CommandPlacements(),
		func(placement domain.CommandPlacement) CommandToCommandGroupModel {
			return CommandToCommandGroupModel{
				CommandGroupId: domainCommandGroup.Id,
				CommandId:      placement.CommandId,
				Position:       placement.Position,
			}
		},
	)
}

// ToDomainCommandGroups folds the rows of commandGroupQuery back into Command
// Groups, keeping the order the query put them in - both the Groups and, within
// each, its Commands.
func ToDomainCommandGroups(rows []commandGroupRow) []domain.CommandGroup {
	commandGroups := make([]domain.CommandGroup, 0)
	indexById := make(map[string]int)

	for _, row := range rows {
		index, alreadyRead := indexById[row.Id]
		if !alreadyRead {
			index = len(commandGroups)
			indexById[row.Id] = index
			commandGroups = append(commandGroups, domain.CommandGroup{
				Id:        row.Id,
				ProjectId: row.ProjectId,
				Name:      row.Name,
				Position:  row.Position,
				Commands:  make([]commanddomain.Command, 0),
			})
		}

		if command, held := row.command(); held {
			commandGroups[index].Commands = append(commandGroups[index].Commands, command)
		}
	}

	return commandGroups
}

func (row commandGroupRow) command() (commanddomain.Command, bool) {
	if row.CommandId == nil {
		return commanddomain.Command{}, false
	}

	return infrastructure.ToDomainCommand(infrastructure.CommandModel{
		Id:               *row.CommandId,
		ProjectId:        row.CommandProjectId,
		Name:             row.CommandName,
		Command:          row.CommandCommand,
		WorkingDirectory: row.CommandWorkingDirectory,
		Position:         row.CommandPosition,
		Link:             row.CommandLink,
		ErrorPatterns:    row.CommandErrorPatterns,
	}), true
}
