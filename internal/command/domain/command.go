package domain

type Command struct {
	Id               string
	ProjectId        string
	Name             string
	Command          string
	WorkingDirectory string
	Position         int
	Link             string
	ErrorPatterns    []string
}
