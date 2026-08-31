package domain

type Repository interface {
	// Get reports a missing Command Group as domainerrors.ErrNotFound, so a nil
	// error guarantees a usable Command Group.
	Get(id string) (CommandGroup, error)
	GetAll(projectId string) ([]CommandGroup, error)
	// GetAllContaining answers the Command Groups holding the Command, each
	// with all of its Commands, across every Project.
	GetAllContaining(commandId string) ([]CommandGroup, error)
	// The three reads below answer what Get, GetAll and GetAllContaining do, in
	// the form that names the Commands instead of carrying them, and report a
	// missing Command Group the same way. They read membership alone, so a
	// Command Group still names a Command whose own record is gone: a rule
	// about membership can be decided without a delete elsewhere having run
	// first.
	GetWithCommandIds(id string) (CommandGroupWithCommandIds, error)
	GetAllWithCommandIds(projectId string) ([]CommandGroupWithCommandIds, error)
	GetAllContainingWithCommandIds(commandId string) ([]CommandGroupWithCommandIds, error)
	Create(commandGroup *CommandGroup) error
	Update(commandGroup *CommandGroup) error
	// UpdateWithCommandIds writes what Update writes, from the form that names
	// the Commands the Command Group holds instead of carrying them.
	UpdateWithCommandIds(commandGroup *CommandGroupWithCommandIds) error
	Delete(commandGroupId string) error
}
