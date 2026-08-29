package execution

import (
	commanddomain "gomander/internal/command/domain"
)

// Runner spawns a Command as a Process in an Environment and answers which
// Commands it still holds one for. It is declared once here, rather than once
// per consumer the way the narrow logging and event roles are, because the
// Environment already travels in its signature: every consumer imports this
// package anyway, and a term the glossary names deserves one definition rather
// than four partial ones.
type Runner interface {
	RunCommand(command *commanddomain.Command, environment Environment) error
	StopRunningCommand(id string) error
	StopAllRunningCommands() []error
	GetRunningCommandIds() []string
}
