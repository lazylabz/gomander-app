package domain

type ConfigRepository interface {
	GetOrCreateConfig() (*Config, error)
	SaveConfig(config *Config) error
}
