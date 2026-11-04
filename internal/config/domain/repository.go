package domain

type Repository interface {
	GetOrCreate() (*Config, error)
	// Update writes everything but the Zoom Level, which only UpdateZoomLevel
	// writes: a zoom overlapping another config write would otherwise undo it,
	// or be undone by it.
	Update(config *Config) error
	UpdateZoomLevel(level ZoomLevel) error
}
