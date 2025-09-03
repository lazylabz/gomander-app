package app

import (
	"errors"

	"gomander/internal/project/domain/event"
)

func (a *App) DeleteProject(projectId string) error {
	err := a.projectRepository.Delete(projectId)
	if err != nil {
		return err
	}

	domainEvent := event.NewProjectDeletedEvent(projectId)

	errs := a.eventBus.PublishSync(domainEvent)

	if len(errs) > 0 {
		combinedErrMsg := "Errors occurred while removing project:"

		for _, pubErr := range errs {
			combinedErrMsg += "\n- " + pubErr.Error()
			a.logger.Error(pubErr.Error())
		}

		err = errors.New(combinedErrMsg)

		return err
	}

	return nil
}
