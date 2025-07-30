package project

import (
	"slices"
)

type CommandGroup struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	CommandIds []string `json:"commands"`
}

func (p *Project) SetCommandGroups(newCommandGroups []CommandGroup) {
	p.CommandGroups = newCommandGroups
}

func (p *Project) GetCommandGroups() []CommandGroup {
	return p.CommandGroups
}

func (p *Project) RemoveCommandFromCommandGroups(commandId string) {
	newCommandGroups := make([]CommandGroup, 0)

	for _, group := range p.CommandGroups {
		if slices.Contains(group.CommandIds, commandId) {
			newCommandIds := make([]string, 0)

			for _, cmdId := range group.CommandIds {
				if cmdId != commandId {
					newCommandIds = append(newCommandIds, cmdId)
				}
			}

			group.CommandIds = newCommandIds
		}
		if len(group.CommandIds) == 0 {
			// If the group has no commands left, we can remove it
			continue
		}
		newCommandGroups = append(newCommandGroups, group)

	}

	p.CommandGroups = newCommandGroups
}
