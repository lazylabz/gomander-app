package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strings"
)

type Event string

type EventEmitter struct {
	ctx context.Context
}

func NewEventEmitter(ctx context.Context) *EventEmitter {
	return &EventEmitter{
		ctx: ctx,
	}
}

const (
	GetCommands         Event = "get_commands"
	ProcessFinished     Event = "process_finished"
	NewLogEntry         Event = "new_log_entry"
	ErrorNotification   Event = "error_notification"
	SuccessNotification Event = "success_notification"
	GetUserConfig       Event = "get_user_config"
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
}

func (e *EventEmitter) emitEvent(event Event, payload interface{}) {
	runtime.EventsEmit(e.ctx, string(event), payload)
}
