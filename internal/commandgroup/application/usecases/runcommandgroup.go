package usecases

import (
	commanddomain "gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/execution"
	"gomander/internal/openedproject"
)

type RunCommandGroup struct {
	openedProject          openedproject.OpenedProject
	commandGroupRepository commandgroupdomain.Repository
	commandRepository      commanddomain.Repository
	commandRunner          execution.Runner
}

func NewRunCommandGroup(
	openedProject openedproject.OpenedProject,
	commandGroupRepo commandgroupdomain.Repository,
	commandRepo commanddomain.Repository,
	runner execution.Runner,
) *RunCommandGroup {
	return &RunCommandGroup{
		openedProject:          openedProject,
		commandGroupRepository: commandGroupRepo,
		commandRepository:      commandRepo,
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

	for _, commandId := range cmdGroup.CommandIds {
		command, err := uc.commandRepository.Get(commandId)
		if err != nil {
			return err
		}

		if err := uc.commandRunner.RunCommand(&command, environment); err != nil {
			return err
		}
	}

	return nil
}
