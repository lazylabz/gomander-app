package domain

import "gomander/internal/ordering"

// Order is where a Command sits among its Project's Commands; nothing else
// assigns a Position.
var Order = ordering.NewList(
	func(command Command) string { return command.Id },
	func(command Command) int { return command.Position },
	func(command *Command, position int) { command.Position = position },
)
