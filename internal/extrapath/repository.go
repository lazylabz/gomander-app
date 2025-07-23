package extrapath

type Repository struct {
	extraPaths []ExtraPath
}

func NewExtraPathRepository(extraPaths []ExtraPath) *Repository {
	return &Repository{
		extraPaths: extraPaths,
	}
}

func (r *Repository) SetExtraPaths(newExtraPaths []ExtraPath) {
	r.extraPaths = newExtraPaths
}

func (r *Repository) GetExtraPaths() []ExtraPath {
	return r.extraPaths
}
