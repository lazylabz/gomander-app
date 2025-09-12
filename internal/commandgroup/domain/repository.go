package domain

type Repository interface {
	Get(id string) (*CommandGroup, error)
	GetAll(projectId string) ([]CommandGroup, error)
	Create(commandGroup *CommandGroup) error
	Update(commandGroup *CommandGroup) error
	Delete(commandGroupId string) error
	RemoveCommandFromCommandGroups(commandId string) error
	DeleteEmpty() (deletedIds []string, err error)
	DeleteAll(projectId string) (deletedIds []string, err error)
}
