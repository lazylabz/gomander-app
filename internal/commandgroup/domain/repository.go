package domain

type Repository interface {
	// Get reports a missing Command Group as domainerrors.ErrNotFound, so a nil
	// error guarantees a usable Command Group.
	Get(id string) (CommandGroup, error)
	GetAll(projectId string) ([]CommandGroup, error)
	Create(commandGroup *CommandGroup) error
	Update(commandGroup *CommandGroup) error
	Delete(commandGroupId string) error
	RemoveCommandFromCommandGroups(commandId string) error
	DeleteEmpty() (deleted []CommandGroup, err error)
	DeleteAll(projectId string) (deletedIds []string, err error)
}
