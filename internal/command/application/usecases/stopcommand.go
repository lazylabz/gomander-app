package usecases

import (
	"gomander/internal/command/domain"
	"gomander/internal/execution"
)

type StopCommand struct {
	commandRepository domain.Repository
	commandRunner     execution.Runner
}

func NewStopCommand(commandRepo domain.Repository, runner execution.Runner) *StopCommand {
	return &StopCommand{
		commandRepository: commandRepo,
		commandRunner:     runner,
	}
}

func (uc *StopCommand) Execute(commandId string) error {
	// A stale UI can ask to stop a Command that is already gone; say so rather
	// than signalling a process nothing points at any more.
	if _, err := uc.commandRepository.Get(commandId); err != nil {
		return err
	}

	return uc.commandRunner.StopRunningCommand(commandId)
}
