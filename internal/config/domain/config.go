package domain

type EnvironmentPath struct {
	Id   string
	Path string
}

type Config struct {
	LastOpenedProjectId string
	EnvironmentPaths    []EnvironmentPath
}
