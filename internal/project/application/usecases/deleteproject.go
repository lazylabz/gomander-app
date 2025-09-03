// filepath: /Users/moises/Code/personal/gomander/internal/project/application/usecases/deleteproject.go
package usecases

import (
	"errors"

	"gomander/internal/eventbus"
	"gomander/internal/logger"
	"gomander/internal/project/domain"
	"gomander/internal/project/domain/event"
)

type DeleteProject interface {
	Execute(projectId string) error
}

type DefaultDeleteProject struct {
	projectRepository domain.Repository
	eventBus          eventbus.EventBus
	logger            logger.Logger
}

func NewDeleteProject(
	projectRepo domain.Repository,
	eventBus eventbus.EventBus,
	logger logger.Logger,
) *DefaultDeleteProject {
	return &DefaultDeleteProject{
		projectRepository: projectRepo,
		eventBus:          eventBus,
		logger:            logger,
	}
}

func (uc *DefaultDeleteProject) Execute(projectId string) error {
	err := uc.projectRepository.Delete(projectId)
	if err != nil {
		return err
	}

	domainEvent := event.NewProjectDeletedEvent(projectId)

	errs := uc.eventBus.PublishSync(domainEvent)

	if len(errs) > 0 {
		combinedErrMsg := "Errors occurred while removing project:"

		for _, pubErr := range errs {
			combinedErrMsg += "\n- " + pubErr.Error()
			uc.logger.Error(pubErr.Error())
		}

		err = errors.New(combinedErrMsg)

		return err
	}

	return nil
}
