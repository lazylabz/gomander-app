package domain

type Repository interface {
	GetAll() ([]Project, error)
	// Get reports a missing Project as domainerrors.ErrNotFound, so a nil error
	// guarantees a usable Project; Find is for the callers where a missing
	// Project is a legitimate outcome.
	Get(id string) (Project, error)
	Find(id string) (project Project, found bool, err error)
	Create(project Project) error
	Update(project Project) error
	Delete(id string) error
}
