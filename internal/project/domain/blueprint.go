package domain

// Blueprint is what a Project looks like before it exists: the name and working
// directory to start from, and the Commands and Command Groups to create. It is
// what reading an importable file produces and what an export describes, so the
// file format that carries it is not its concern.
type Blueprint struct {
	Name             string
	WorkingDirectory string
	Commands         []BlueprintCommand
	CommandGroups    []BlueprintCommandGroup
}

// BlueprintCommand carries an Id that only means something inside the
// Blueprint: it is how a BlueprintCommandGroup names its members. The Command
// the import creates gets an Id of its own.
type BlueprintCommand struct {
	Id               string
	Name             string
	Command          string
	WorkingDirectory string
}

type BlueprintCommandGroup struct {
	Id         string
	Name       string
	CommandIds []string
}
