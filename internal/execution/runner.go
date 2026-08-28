package execution

import (
	commanddomain "gomander/internal/command/domain"
)

// Runner spawns a Command as a Process in an Environment and answers which
// Commands it still holds one for. It is declared here rather than next to the
// process runner because every consumer of it already speaks in Environments.
type Runner interface {
	RunCommand(command *commanddomain.Command, environment Environment) error
	StopRunningCommand(id string) error
	StopAllRunningCommands() []error
	GetRunningCommandIds() []string
}
