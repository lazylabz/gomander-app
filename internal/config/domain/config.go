package domain

type EnvironmentPath struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

type Config struct {
	LastOpenedProjectId string            `json:"lastOpenedProjectId"`
	EnvironmentPaths    []EnvironmentPath `json:"environmentPaths"`
}
