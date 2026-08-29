package usecases

import (
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/execution"
	"gomander/internal/openedproject"
)

type RunCommandGroup struct {
	openedProject          openedproject.OpenedProject
	commandGroupRepository commandgroupdomain.Repository
	commandRunner          execution.Runner
}

func NewRunCommandGroup(
	openedProject openedproject.OpenedProject,
	commandGroupRepo commandgroupdomain.Repository,
	runner execution.Runner,
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
