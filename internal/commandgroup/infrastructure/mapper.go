package infrastructure

import (
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
	return array.Map(domainCommandGroup.CommandPlacements(), func(placement domain.CommandPlacement) CommandToCommandGroupModel {
		return CommandToCommandGroupModel{
			CommandGroupId: domainCommandGroup.Id,
			CommandId:      placement.CommandId,
			Position:       placement.Position,
		}
	})
}

// ToDomainCommandGroups folds the rows of commandGroupIdentityQuery back into
// Command Groups, keeping the order the query put them in - both the Groups
// and, within each, the Commands it names.
func ToDomainCommandGroups(rows []commandGroupIdentityRow) []domain.CommandGroup {
	commandGroups := make([]domain.CommandGroup, 0)
	indexById := make(map[string]int)

	for _, row := range rows {
		index, alreadyRead := indexById[row.Id]
		if !alreadyRead {
			index = len(commandGroups)
			indexById[row.Id] = index
			commandGroups = append(commandGroups, domain.CommandGroup{
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
