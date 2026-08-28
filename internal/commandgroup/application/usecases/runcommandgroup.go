package usecases

import (
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/openedproject"
	"gomander/internal/runner"
)

type RunCommandGroup struct {
	openedProject          openedproject.OpenedProject
	commandGroupRepository commandgroupdomain.Repository
	commandRunner          runner.Runner
}

func NewRunCommandGroup(
	openedProject openedproject.OpenedProject,
	commandGroupRepo commandgroupdomain.Repository,
	runner runner.Runner,
) *RunCommandGroup {
	return &RunCommandGroup{
		openedProject:          openedProject,
		commandGroupRepository: commandGroupRepo,
		commandRunner:          runner,
	}
}

func (uc *RunCommandGroup) Execute(commandGroupId string) error {
	cmdGroup, err := uc.commandGroupRepository.Get(commandGroupId)
	if err != nil {
		return err
	}

	environment, err := uc.openedProject.ExecutionEnvironment()
	if err != nil {
		return err
	}

	for _, command := range cmdGroup.Commands {
		if err := uc.commandRunner.RunCommand(&command, environment); err != nil {
			return err
		}
	}

	return nil
}
