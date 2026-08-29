package usecases

import (
	"gomander/internal/eventbus"
	"gomander/internal/project/domain"
	"gomander/internal/project/domain/event"
)

type DeleteProject struct {
	projectRepository domain.Repository
	eventBus          eventbus.EventBus
}

func NewDeleteProject(
	projectRepo domain.Repository,
	eventBus eventbus.EventBus,
) *DeleteProject {
	return &DeleteProject{
		projectRepository: projectRepo,
		eventBus:          eventBus,
	}
}

func (uc *DeleteProject) Execute(projectId string) error {
	err := uc.projectRepository.Delete(projectId)
	if err != nil {
		return err
	}

	return eventbus.Combined(
		"Errors occurred while removing project:",
		uc.eventBus.PublishSync(event.NewProjectDeletedEvent(projectId)),
	)
}
