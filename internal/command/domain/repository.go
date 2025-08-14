package domain

type Repository interface {
	Get(commandId string) (*Command, error)
	GetAll(projectId string) ([]Command, error)
	Create(command *Command) error
	Update(command *Command) error
	Delete(commandId string) error
}
