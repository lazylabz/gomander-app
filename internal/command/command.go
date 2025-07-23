package command

type Command struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Command          string `json:"command"`
	WorkingDirectory string `json:"workingDirectory"`
}
