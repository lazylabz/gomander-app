package domain

import (
	"gomander/internal/helpers/array"
)

// Cascade is what a set of Command Groups becomes once a Command they held is
// gone: the ones still standing, and the ones that ceased to exist because that
// Command was their last.
type Cascade struct {
	Survived []CommandGroupWithCommandIds
	Deleted  []CommandGroupWithCommandIds
}

// RemoveCommandFrom drops the Command from every Command Group given and
// applies the rule that a Command Group with no Commands ceases to exist. It
// answers only what is left; persisting that answer atomically is the caller's.
//
// Give it the Command Groups that held the Command, which
// Repository.GetAllContainingWithCommandIds answers from the membership table.
// The rule is decided on membership alone, so it answers the same whether or
// not the Command's own record has already been deleted.
func RemoveCommandFrom(commandGroups []CommandGroupWithCommandIds, commandId string) Cascade {
	cascade := Cascade{
		Survived: make([]CommandGroupWithCommandIds, 0, len(commandGroups)),
		Deleted:  make([]CommandGroupWithCommandIds, 0),
	}

	for _, commandGroup := range commandGroups {
		commandGroup.CommandIds = array.Filter(commandGroup.CommandIds, func(heldCommandId string) bool {
			return heldCommandId != commandId
		})

		if len(commandGroup.CommandIds) == 0 {
			cascade.Deleted = append(cascade.Deleted, commandGroup)
			continue
		}

		cascade.Survived = append(cascade.Survived, commandGroup)
	}

	return cascade
}
