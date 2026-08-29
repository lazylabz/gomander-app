package domain

import (
	commanddomain "gomander/internal/command/domain"
	"gomander/internal/helpers/array"
)

// Cascade is what a set of Command Groups becomes once a Command they held is
// gone: the ones still standing, and the ones that ceased to exist because that
// Command was their last.
type Cascade struct {
	Survived []CommandGroup
	Deleted  []CommandGroup
}

// RemoveCommandFrom drops the Command from every Command Group given and
// applies the rule that a Command Group with no Commands ceases to exist. It
// answers only what is left; persisting that answer atomically is the caller's.
//
// Give it the Command Groups that held the Command, which
// Repository.GetAllContaining answers from the join table. The Command row is
// already deleted by the time this runs, so a Group arrives carrying only its
// surviving Commands and the filter usually finds nothing to drop: emptiness
// here means the Command was its last, and a Group empty for any other reason
// must not be in the input.
func RemoveCommandFrom(commandGroups []CommandGroup, commandId string) Cascade {
	cascade := Cascade{
		Survived: make([]CommandGroup, 0, len(commandGroups)),
		Deleted:  make([]CommandGroup, 0),
	}

	for _, commandGroup := range commandGroups {
		commandGroup.Commands = array.Filter(commandGroup.Commands, func(command commanddomain.Command) bool {
			return command.Id != commandId
		})

		if len(commandGroup.Commands) == 0 {
			cascade.Deleted = append(cascade.Deleted, commandGroup)
			continue
		}

		cascade.Survived = append(cascade.Survived, commandGroup)
	}

	return cascade
}
