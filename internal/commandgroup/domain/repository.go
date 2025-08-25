package domain

type Repository interface {
	Get(id string) (*CommandGroup, error)
	GetAll(projectId string) ([]CommandGroup, error)
	Create(commandGroup *CommandGroup) error
	Update(commandGroup *CommandGroup) error
	Delete(commandGroupId string) error
	RemoveCommandFromCommandGroups(commandId string) error
	DeleteEmpty() error
}
