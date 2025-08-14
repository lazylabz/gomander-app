package domain

type Repository interface {
	GetOrCreate() (*Config, error)
	Update(config *Config) error
}
