package usecases

import (
	"gomander/internal/eventbus"
	"gomander/internal/logger"
	"gomander/internal/project/domain"
	"gomander/internal/project/domain/event"
)

type DeleteProject struct {
	projectRepository domain.Repository
	eventBus          eventbus.EventBus
	logger            logger.Logger
}

func NewDeleteProject(
	projectRepo domain.Repository,
	eventBus eventbus.EventBus,
	logger logger.Logger,
) *DeleteProject {
	return &DeleteProject{
		projectRepository: projectRepo,
		eventBus:          eventBus,
		logger:            logger,
	}
}

func (uc *DeleteProject) Execute(projectId string) error {
	err := uc.projectRepository.Delete(projectId)
	if err != nil {
		return err
	}

	domainEvent := event.NewProjectDeletedEvent(projectId)

	errs := uc.eventBus.PublishSync(domainEvent)

	for _, pubErr := range errs {
		uc.logger.Error(pubErr.Error())
	}

	return eventbus.Combined("Errors occurred while removing project:", errs)
}
