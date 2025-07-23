package commandgroup

import "slices"

type Repository struct {
	commandGroups []CommandGroup
}

func NewCommandGroupRepository(commandGroups []CommandGroup) *Repository {
	return &Repository{
		commandGroups: commandGroups,
	}
}

func (r *Repository) SetCommandGroups(newCommandGroups []CommandGroup) {
	r.commandGroups = newCommandGroups
}

func (r *Repository) GetCommandGroups() []CommandGroup {
	return r.commandGroups
}

func (r *Repository) RemoveCommandFromCommandGroups(commandId string) {
	newCommandGroups := make([]CommandGroup, 0)

	for _, group := range r.commandGroups {
		if slices.Contains(group.CommandIds, commandId) {
			newCommandIds := make([]string, 0)

			for _, cmdId := range group.CommandIds {
				if cmdId != commandId {
					newCommandIds = append(newCommandIds, cmdId)
				}
			}

			group.CommandIds = newCommandIds
		}
		newCommandGroups = append(newCommandGroups, group)

	}

	r.commandGroups = newCommandGroups
}
