package domain

// Status is what a Command is doing right now, as far as the Runner is
// concerned.
type Status int

const (
	Stopped Status = iota
	Running
)

// RunningCommands is the set of Commands the Runner currently holds a Process
// for. It is the single owner of "this Command is running": membership, a
// Command's Status and how many of a Command Group's Commands are running all
// come from here, so no two consumers can answer the question differently.
type RunningCommands struct {
	ids map[string]struct{}
}

func NewRunningCommands(commandIds []string) RunningCommands {
	ids := make(map[string]struct{}, len(commandIds))

	for _, id := range commandIds {
		ids[id] = struct{}{}
	}

	return RunningCommands{ids: ids}
}

func (r RunningCommands) IsRunning(commandId string) bool {
	_, running := r.ids[commandId]
	return running
}

func (r RunningCommands) StatusOf(commandId string) Status {
	if r.IsRunning(commandId) {
		return Running
	}

	return Stopped
}

// CountIn answers how many of the named Commands are running - the Commands a
// Command Group names, so its count and each Command's Status come from the
// same place.
func (r RunningCommands) CountIn(commandIds []string) int {
	count := 0

	for _, commandId := range commandIds {
		if r.IsRunning(commandId) {
			count++
		}
	}

	return count
}
