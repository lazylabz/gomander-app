package domain

type Repository interface {
	// Get, GetAll, GetAllContaining, Create and Update carry whole Commands.
	// Only GetAll still has a caller - the read path - and all five go once it
	// names the Commands it answers instead of carrying them.
	//
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
	// The two writes below write what Create and Update write, from the form
	// that names the Commands the Command Group holds instead of carrying them.
	CreateWithCommandIds(commandGroup *CommandGroupWithCommandIds) error
	UpdateWithCommandIds(commandGroup *CommandGroupWithCommandIds) error
	Delete(commandGroupId string) error
}
