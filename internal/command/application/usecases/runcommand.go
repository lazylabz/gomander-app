package usecases

import (
	"gomander/internal/command/domain"
	"gomander/internal/openedproject"
	"gomander/internal/runner"
)

type RunCommand struct {
	openedProject     openedproject.OpenedProject
	commandRepository domain.Repository
	commandRunner     runner.Runner
}

func NewRunCommand(
	openedProject openedproject.OpenedProject,
	commandRepo domain.Repository,
	runner runner.Runner,
) *RunCommand {
	return &RunCommand{
		openedProject:     openedProject,
		commandRepository: commandRepo,
		commandRunner:     runner,
	}
}

func (uc *RunCommand) Execute(commandId string) error {
	cmd, err := uc.commandRepository.Get(commandId)
	if err != nil {
		return err
	}

	environment, err := uc.openedProject.ExecutionEnvironment()
	if err != nil {
		return err
	}

	return uc.commandRunner.RunCommand(&cmd, environment)
}
