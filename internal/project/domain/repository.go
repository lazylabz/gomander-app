package domain

type Repository interface {
	GetAll() ([]Project, error)
	Get(id string) (*Project, error)
	Create(project Project) error
	Update(project Project) error
	Delete(id string) error
}
