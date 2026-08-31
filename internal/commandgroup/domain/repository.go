package domain

type Repository interface {
	// The three reads below answer from membership alone, so a Command Group
	// still names a Command whose own record is gone: a rule about membership
	// can be decided without a delete elsewhere having run first.
	//
	// Get reports a missing Command Group as domainerrors.ErrNotFound, so a nil
	// error guarantees a usable Command Group.
	Get(id string) (CommandGroup, error)
	GetAll(projectId string) ([]CommandGroup, error)
	// GetAllContaining answers the Command Groups holding the Command, across
	// every Project.
	GetAllContaining(commandId string) ([]CommandGroup, error)
	Create(commandGroup *CommandGroup) error
	Update(commandGroup *CommandGroup) error
	Delete(commandGroupId string) error
}
