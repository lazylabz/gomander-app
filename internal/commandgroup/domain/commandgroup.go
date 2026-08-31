package domain

// CommandGroup names the Commands it holds, in the order it holds them.
type CommandGroup struct {
	Id         string
	ProjectId  string
	Name       string
	CommandIds []string
	Position   int
}
