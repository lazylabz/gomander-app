package domain

import "gomander/internal/ordering"

// Order is where a Command Group sits among its Project's Command Groups;
// nothing else assigns a Position.
var Order = ordering.NewList(
	func(commandGroup CommandGroup) string { return commandGroup.Id },
	func(commandGroup CommandGroup) int { return commandGroup.Position },
	func(commandGroup *CommandGroup, position int) { commandGroup.Position = position },
)
