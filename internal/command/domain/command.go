package domain

type Command struct {
	Id               string `json:"id"`
	ProjectId        string `json:"projectId"`
	Name             string `json:"name"`
	Command          string `json:"command"`
	WorkingDirectory string `json:"workingDirectory"`
	Position         int    `json:"position"`
	Link             string `json:"link"`
	ErrorPatterns    string `json:"errorPatterns"`
}
