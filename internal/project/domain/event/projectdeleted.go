package event

type ProjectDeletedEvent struct {
	ProjectId string
}

func (ProjectDeletedEvent) GetName() string {
	return "domain_event.project.delete"
}

func NewProjectDeletedEvent(projectId string) ProjectDeletedEvent {
	return ProjectDeletedEvent{
		ProjectId: projectId,
	}
}
