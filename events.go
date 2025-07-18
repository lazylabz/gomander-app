package main

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strings"
)

type Event string

const (
	GetCommands       Event = "get_commands"
	ProcessFinished   Event = "process_finished"
	NewLogEntry       Event = "new_log_entry"
	ErrorNotification Event = "error_notification"
)

var Events = []struct {
	Value  Event
	TSName string
}{
	{Value: GetCommands, TSName: strings.ToUpper(string(GetCommands))},
	{Value: ProcessFinished, TSName: strings.ToUpper(string(ProcessFinished))},
	{Value: NewLogEntry, TSName: strings.ToUpper(string(NewLogEntry))},
	{Value: ErrorNotification, TSName: strings.ToUpper(string(ErrorNotification))},
}

func (a *App) emitEvent(event Event, payload interface{}) {
	runtime.EventsEmit(a.ctx, string(event), payload)
}
