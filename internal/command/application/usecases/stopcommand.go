package usecases

import (
	"gomander/internal/command/domain"
	"gomander/internal/runner"
)

type StopCommand interface {
	Execute(commandId string) error
}

type DefaultStopCommand struct {
	commandRepository domain.Repository
	commandRunner     runner.Runner
}

func NewStopCommand(commandRepo domain.Repository, runner runner.Runner) *DefaultStopCommand {
	return &DefaultStopCommand{
		commandRepository: commandRepo,
		commandRunner:     runner,
	}
}

func (uc *DefaultStopCommand) Execute(commandId string) error {
	// A stale UI can ask to stop a Command that is already gone; say so rather
	// than signalling a process nothing points at any more.
	if _, err := uc.commandRepository.Get(commandId); err != nil {
		return err
	}

	return uc.commandRunner.StopRunningCommand(commandId)
}
