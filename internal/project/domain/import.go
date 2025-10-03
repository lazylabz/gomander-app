package domain

type CommandGroupJSONv1 struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	CommandIds []string `json:"commandIds"`
}

type CommandJSONv1 struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Command          string `json:"command"`
	WorkingDirectory string `json:"workingDirectory"`
}

type ProjectExportJSONv1 struct {
	Version          int                  `json:"version"`
	Name             string               `json:"name"`
	WorkingDirectory string               `json:"workingDirectory"`
	Commands         []CommandJSONv1      `json:"commands"`
	CommandGroups    []CommandGroupJSONv1 `json:"commandGroups"`
}
