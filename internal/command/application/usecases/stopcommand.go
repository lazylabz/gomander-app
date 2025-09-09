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
	// Check if the command exists before trying to stop it
	_, err := uc.commandRepository.Get(commandId)
	if err != nil {
		return err
	}

	err = uc.commandRunner.StopRunningCommand(commandId)
	if err != nil {
		return err
	}

	return nil
}
