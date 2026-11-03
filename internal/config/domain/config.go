package domain

type EnvironmentPath struct {
	Id   string
	Path string
}

type Config struct {
	LastOpenedProjectId string
	EnvironmentPaths    []EnvironmentPath
	Locale              string
	LogsFontSize        int
	ZoomLevel           ZoomLevel
}

func (c Config) Validate() error {
	return validateLogsFontSize(c.LogsFontSize)
}
