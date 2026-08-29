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
	Delete(commandGroupId string) error
	// Atomically runs change against a Repository bound to a single
	// transaction: every write inside it lands, or none of them does. It is how
	// a rule the domain decides over several Command Groups reaches storage in
	// one piece.
	Atomically(change func(Repository) error) error
}
