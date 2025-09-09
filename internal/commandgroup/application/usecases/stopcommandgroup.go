package usecases

import (
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/runner"
)

type StopCommandGroup interface {
	Execute(commandGroupId string) error
}

type DefaultStopCommandGroup struct {
	commandGroupRepository commandgroupdomain.Repository
	commandRunner          runner.Runner
}

func NewStopCommandGroup(
	commandGroupRepo commandgroupdomain.Repository,
	runner runner.Runner,
) *DefaultStopCommandGroup {
	return &DefaultStopCommandGroup{
		commandGroupRepository: commandGroupRepo,
		commandRunner:          runner,
	}
}

func (uc *DefaultStopCommandGroup) Execute(commandGroupId string) error {
	cmdGroup, err := uc.commandGroupRepository.Get(commandGroupId)
	if err != nil {
		return err
	}

	err = uc.commandRunner.StopRunningCommands(cmdGroup.Commands)
	if err != nil {
		return err
	}

	return nil
}
