package event

import (
	"strings"
)

type Event string

const (
	GetCommands         Event = "get_commands"
	ProcessFinished     Event = "process_finished"
	NewLogEntry         Event = "new_log_entry"
	ErrorNotification   Event = "error_notification"
	SuccessNotification Event = "success_notification"
	GetUserConfig       Event = "get_user_config"
	GetCommandGroups    Event = "get_command_groups"
)

var Events = []struct {
	Value  Event
	TSName string
}{
	{Value: GetCommands, TSName: strings.ToUpper(string(GetCommands))},
	{Value: ProcessFinished, TSName: strings.ToUpper(string(ProcessFinished))},
	{Value: NewLogEntry, TSName: strings.ToUpper(string(NewLogEntry))},
	{Value: ErrorNotification, TSName: strings.ToUpper(string(ErrorNotification))},
	{Value: SuccessNotification, TSName: strings.ToUpper(string(SuccessNotification))},
	{Value: GetUserConfig, TSName: strings.ToUpper(string(GetUserConfig))},
	{Value: GetCommandGroups, TSName: strings.ToUpper(string(GetCommandGroups))},
}
