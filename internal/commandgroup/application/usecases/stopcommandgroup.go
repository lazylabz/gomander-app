package usecases

import (
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/runner"
)

type StopCommandGroup struct {
	commandGroupRepository commandgroupdomain.Repository
	commandRunner          runner.Runner
}

func NewStopCommandGroup(
	commandGroupRepo commandgroupdomain.Repository,
	runner runner.Runner,
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

	for _, command := range cmdGroup.Commands {
		if err := uc.commandRunner.StopRunningCommand(command.Id); err != nil {
			return err
		}
	}

	return nil
}
