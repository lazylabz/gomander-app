package domain

type Repository interface {
	// Get reports a missing Command Group as domainerrors.ErrNotFound, so a nil
	// error guarantees a usable Command Group.
	Get(id string) (CommandGroup, error)
	GetAll(projectId string) ([]CommandGroup, error)
	// GetAllContaining answers the Command Groups holding the Command, each
	// with all of its Commands, across every Project.
	GetAllContaining(commandId string) ([]CommandGroup, error)
	Create(commandGroup *CommandGroup) error
	Update(commandGroup *CommandGroup) error
	Delete(commandGroupId string) error
	// Atomically runs change against a Repository bound to a single
	// transaction: every write inside it lands, or none of them does. It is how
	// a rule the domain decides over several Command Groups reaches storage in
	// one piece.
	Atomically(change func(Repository) error) error
}
