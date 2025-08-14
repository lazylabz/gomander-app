package domain

type Repository interface {
	GetCommandGroups(projectId string) ([]CommandGroup, error)
	GetCommandGroupById(id string) (*CommandGroup, error)
	CreateCommandGroup(commandGroup *CommandGroup) error
	UpdateCommandGroup(commandGroup *CommandGroup) error
	DeleteCommandGroup(commandGroupId string) error
	UpdateCommandGroupCommands(group *CommandGroup) error
}
