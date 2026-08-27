package domain

type Repository interface {
	// Get reports a missing Command as domainerrors.ErrNotFound, so a nil error
	// guarantees a usable Command.
	Get(commandId string) (Command, error)
	GetAll(projectId string) ([]Command, error)
	Create(command *Command) error
	Update(command *Command) error
	Delete(commandId string) error
	DeleteAll(projectId string) error
}
