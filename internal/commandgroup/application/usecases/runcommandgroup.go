package usecases

import (
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/openedproject"
	"gomander/internal/runner"
)

type RunCommandGroup interface {
	Execute(commandGroupId string) error
}

type DefaultRunCommandGroup struct {
	openedProject          openedproject.OpenedProject
	commandGroupRepository commandgroupdomain.Repository
	commandRunner          runner.Runner
}

func NewRunCommandGroup(
	openedProject openedproject.OpenedProject,
	commandGroupRepo commandgroupdomain.Repository,
	runner runner.Runner,
) *DefaultRunCommandGroup {
	return &DefaultRunCommandGroup{
		openedProject:          openedProject,
		commandGroupRepository: commandGroupRepo,
		commandRunner:          runner,
	}
}

func (uc *DefaultRunCommandGroup) Execute(commandGroupId string) error {
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
