package config

type UserConfig struct {
	EnvironmentPaths    []string `json:"environmentPaths"`
	LastOpenedProjectId string   `json:"lastOpenedProjectId"`
}

func EmptyUserConfig() *UserConfig {
	return &UserConfig{
		EnvironmentPaths:    make([]string, 0),
		LastOpenedProjectId: "",
	}
}
