package usecases

import (
	"errors"

	"github.com/google/uuid"

	"gomander/internal/command/domain"
	domainevent "gomander/internal/command/domain/event"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/eventbus"
)

type DuplicateCommand interface {
	Execute(commandId, targetGroupId string) error
}

type DefaultDuplicateCommand struct {
	configRepository  configdomain.Repository
	commandRepository domain.Repository
	eventBus          eventbus.EventBus
}

func NewDuplicateCommand(configRepo configdomain.Repository, commandRepo domain.Repository, eventBus eventbus.EventBus) *DefaultDuplicateCommand {
	return &DefaultDuplicateCommand{
		configRepository:  configRepo,
		commandRepository: commandRepo,
		eventBus:          eventBus,
	}
}

func (uc *DefaultDuplicateCommand) Execute(commandId, targetGroupId string) error {
	userConfig, err := uc.configRepository.GetOrCreate()
	if err != nil {
		return err
	}

	originalCommand, err := uc.commandRepository.Get(commandId)
	if err != nil {
		return err
	}

	allCommands, err := uc.commandRepository.GetAll(userConfig.LastOpenedProjectId)
	if err != nil {
		return err
	}

	duplicatedCommand := domain.Command{
		Id:               uuid.New().String(),
		ProjectId:        originalCommand.ProjectId,
		Name:             originalCommand.Name + " (copy)",
		Command:          originalCommand.Command,
		WorkingDirectory: originalCommand.WorkingDirectory,
		Position:         len(allCommands),
	}

	err = uc.commandRepository.Create(&duplicatedCommand)
	if err != nil {
		return err
	}

	domainEvent := domainevent.NewCommandDuplicatedEvent(duplicatedCommand.Id, targetGroupId)

	errs := uc.eventBus.PublishSync(domainEvent)

	if len(errs) > 0 {
		combinedErrMsg := "Errors occurred while duplicating command:"

		for _, pubErr := range errs {
			combinedErrMsg += "\n- " + pubErr.Error()
		}

		err = errors.New(combinedErrMsg)

		return err
	}

	return nil
}
