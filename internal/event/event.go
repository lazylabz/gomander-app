package event

import (
	"strings"
)

type Event string

const (
	ProcessStarted  Event = "process_started"
	ProcessFinished Event = "process_finished"
	NewLogEntry     Event = "new_log_entry"
	GetUserConfig   Event = "get_user_config"
)

var Events = []struct {
	Value  Event
	TSName string
}{
	{Value: ProcessFinished, TSName: strings.ToUpper(string(ProcessFinished))},
	{Value: ProcessStarted, TSName: strings.ToUpper(string(ProcessStarted))},
	{Value: NewLogEntry, TSName: strings.ToUpper(string(NewLogEntry))},
	{Value: GetUserConfig, TSName: strings.ToUpper(string(GetUserConfig))},
}
