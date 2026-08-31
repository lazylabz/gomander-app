package usecases

import (
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/execution"
)

type StopCommandGroup struct {
	commandGroupRepository commandgroupdomain.Repository
	commandRunner          execution.Runner
}

func NewStopCommandGroup(
	commandGroupRepo commandgroupdomain.Repository,
	runner execution.Runner,
) *StopCommandGroup {
	return &StopCommandGroup{
		commandGroupRepository: commandGroupRepo,
		commandRunner:          runner,
	}
}

func (uc *StopCommandGroup) Execute(commandGroupId string) error {
	cmdGroup, err := uc.commandGroupRepository.Get(commandGroupId)
	if err != nil {
		return err
	}

	for _, commandId := range cmdGroup.CommandIds {
		if err := uc.commandRunner.StopRunningCommand(commandId); err != nil {
			return err
		}
	}

	return nil
}
