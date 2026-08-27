package usecases

import (
	"github.com/google/uuid"

	"gomander/internal/command/domain"
	domainevent "gomander/internal/command/domain/event"
	"gomander/internal/eventbus"
	"gomander/internal/openedproject"
)

type DuplicateCommand struct {
	openedProject     openedproject.OpenedProject
	commandRepository domain.Repository
	eventBus          eventbus.EventBus
}

func NewDuplicateCommand(openedProject openedproject.OpenedProject, commandRepo domain.Repository, eventBus eventbus.EventBus) *DuplicateCommand {
	return &DuplicateCommand{
		openedProject:     openedProject,
		commandRepository: commandRepo,
		eventBus:          eventBus,
	}
}

func (uc *DuplicateCommand) Execute(commandId, targetGroupId string) error {
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
	duplicatedCommand.Position = domain.Order.End(allCommands)

	err = uc.commandRepository.Create(&duplicatedCommand)
	if err != nil {
		return err
	}

	domainEvent := domainevent.NewCommandDuplicatedEvent(duplicatedCommand.Id, targetGroupId)

	return eventbus.Combined(
		"Errors occurred while duplicating command:",
		uc.eventBus.PublishSync(domainEvent),
	)
}
