package domain

import "gomander/internal/ordering"

// Order is where a Command Group sits among its Project's Command Groups;
// nothing else assigns a Position.
var Order = ordering.NewList(
	func(commandGroup CommandGroupWithCommandIds) string { return commandGroup.Id },
	func(commandGroup CommandGroupWithCommandIds) int { return commandGroup.Position },
	func(commandGroup *CommandGroupWithCommandIds, position int) { commandGroup.Position = position },
)

// CommandPlacement is where one of the Commands a Command Group holds sits
// among the others.
type CommandPlacement struct {
	CommandId string
	Position  int
}

// CommandOrder is where a Command sits among the Commands of the Command Group
// holding it; nothing else assigns a Position.
var CommandOrder = ordering.NewList(
	func(placement CommandPlacement) string { return placement.CommandId },
	func(placement CommandPlacement) int { return placement.Position },
	func(placement *CommandPlacement, position int) { placement.Position = position },
)

// CommandPlacements answers where each Command the Command Group holds sits: in
// the order the Group holds them, each one placed behind the last. Storage
// writes these Positions rather than deciding any of its own.
func (commandGroup CommandGroup) CommandPlacements() []CommandPlacement {
	placements := make([]CommandPlacement, 0, len(commandGroup.Commands))

	for _, command := range commandGroup.Commands {
		placements = append(placements, CommandPlacement{
			CommandId: command.Id,
			Position:  CommandOrder.End(placements),
		})
	}

	return placements
}

// CommandPlacements answers the same for a Command Group that names the
// Commands it holds instead of carrying them.
func (commandGroup CommandGroupWithCommandIds) CommandPlacements() []CommandPlacement {
	placements := make([]CommandPlacement, 0, len(commandGroup.CommandIds))

	for _, commandId := range commandGroup.CommandIds {
		placements = append(placements, CommandPlacement{
			CommandId: commandId,
			Position:  CommandOrder.End(placements),
		})
	}

	return placements
}
