package event

type Event string

const (
	ProcessStarted       Event = "process_started"
	ProcessFinished      Event = "process_finished"
	NewLogEntry          Event = "new_log_entry"
	CommandGroupDeleted  Event = "command_group_deleted"
	CommandErrorDetected Event = "command_error_detected"
)

// LogEntryKind travels with every NewLogEntry so the layer that displays a line
// can tell the Command's own text from the Process output it produced. How
// either one looks is decided there, not here.
type LogEntryKind string

const (
	CommandLogEntry LogEntryKind = "command"
	OutputLogEntry  LogEntryKind = "output"
)
