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

func ToCommandGroupModelWithCommandIds(domainCommandGroup *domain.CommandGroupWithCommandIds) CommandGroupModel {
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
	return toCommandToCommandGroupModels(domainCommandGroup.Id, domainCommandGroup.CommandPlacements())
}

func ToCommandToCommandGroupModelsWithCommandIds(domainCommandGroup *domain.CommandGroupWithCommandIds) []CommandToCommandGroupModel {
	return toCommandToCommandGroupModels(domainCommandGroup.Id, domainCommandGroup.CommandPlacements())
}

func toCommandToCommandGroupModels(commandGroupId string, placements []domain.CommandPlacement) []CommandToCommandGroupModel {
	return array.Map(placements, func(placement domain.CommandPlacement) CommandToCommandGroupModel {
		return CommandToCommandGroupModel{
			CommandGroupId: commandGroupId,
			CommandId:      placement.CommandId,
			Position:       placement.Position,
		}
	})
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

// ToDomainCommandGroupsWithCommandIds folds the rows of
// commandGroupIdentityQuery back into Command Groups, keeping the order the
// query put them in - both the Groups and, within each, the Commands it names.
func ToDomainCommandGroupsWithCommandIds(rows []commandGroupIdentityRow) []domain.CommandGroupWithCommandIds {
	commandGroups := make([]domain.CommandGroupWithCommandIds, 0)
	indexById := make(map[string]int)

	for _, row := range rows {
		index, alreadyRead := indexById[row.Id]
		if !alreadyRead {
			index = len(commandGroups)
			indexById[row.Id] = index
			commandGroups = append(commandGroups, domain.CommandGroupWithCommandIds{
				Id:         row.Id,
				ProjectId:  row.ProjectId,
				Name:       row.Name,
				Position:   row.Position,
				CommandIds: make([]string, 0),
			})
		}

		if row.CommandId != nil {
			commandGroups[index].CommandIds = append(commandGroups[index].CommandIds, *row.CommandId)
		}
	}

	return commandGroups
}
