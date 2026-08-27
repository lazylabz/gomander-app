package infrastructure

import (
	commanddomain "gomander/internal/command/domain"
	"gomander/internal/command/infrastructure"
	"gomander/internal/commandgroup/domain"
)

func ToCommandGroupModel(domainCommandGroup *domain.CommandGroup) CommandGroupModel {
	return CommandGroupModel{
		Id:        domainCommandGroup.Id,
		Name:      domainCommandGroup.Name,
		ProjectId: domainCommandGroup.ProjectId,
		Position:  domainCommandGroup.Position,
	}
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
