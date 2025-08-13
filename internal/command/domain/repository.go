package domain

type Repository interface {
	GetCommand(commandId string) (*Command, error)
	GetCommands(projectId string) ([]*Command, error)
	SaveCommand(command *Command) error
	EditCommand(command *Command) error
	DeleteCommand(commandId string) error
}
