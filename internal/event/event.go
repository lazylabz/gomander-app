package event

import (
	"strings"
)

type Event string

const (
	ProcessStarted      Event = "process_started"
	ProcessFinished     Event = "process_finished"
	NewLogEntry         Event = "new_log_entry"
	CommandGroupDeleted Event = "command_group_deleted"
	CommandFailed       Event = "command_failed"
)

var Events = []struct {
	Value  Event
	TSName string
}{
	{Value: ProcessFinished, TSName: strings.ToUpper(string(ProcessFinished))},
	{Value: ProcessStarted, TSName: strings.ToUpper(string(ProcessStarted))},
	{Value: NewLogEntry, TSName: strings.ToUpper(string(NewLogEntry))},
	{Value: CommandGroupDeleted, TSName: strings.ToUpper(string(CommandGroupDeleted))},
	{Value: CommandFailed, TSName: strings.ToUpper(string(CommandFailed))},
}
