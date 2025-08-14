package domain

type Repository interface {
	GetOrCreateConfig() (*Config, error)
	SaveConfig(config *Config) error
}
