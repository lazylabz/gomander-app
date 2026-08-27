package usecases

import (
	"errors"

	"github.com/google/uuid"

	"gomander/internal/command/domain"
	domainevent "gomander/internal/command/domain/event"
	"gomander/internal/eventbus"
	"gomander/internal/openedproject"
)

type DuplicateCommand interface {
	Execute(commandId, targetGroupId string) error
}

type DefaultDuplicateCommand struct {
	openedProject     openedproject.OpenedProject
	commandRepository domain.Repository
	eventBus          eventbus.EventBus
}

func NewDuplicateCommand(openedProject openedproject.OpenedProject, commandRepo domain.Repository, eventBus eventbus.EventBus) *DefaultDuplicateCommand {
	return &DefaultDuplicateCommand{
		openedProject:     openedProject,
		commandRepository: commandRepo,
		eventBus:          eventBus,
	}
}

func (uc *DefaultDuplicateCommand) Execute(commandId, targetGroupId string) error {
	project, err := uc.openedProject.Get()
	if err != nil {
		return err
	}

	originalCommand, err := uc.commandRepository.Get(commandId)
	if err != nil {
		return err
	}

	allCommands, err := uc.commandRepository.GetAll(project.Id)
	if err != nil {
		return err
	}

	duplicatedCommand := *originalCommand

	// Override specific fields
	duplicatedCommand.Id = uuid.New().String()
	duplicatedCommand.Name = originalCommand.Name + " (copy)"
	duplicatedCommand.Position = len(allCommands)

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
