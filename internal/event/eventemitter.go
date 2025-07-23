package event

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type EventEmitter struct {
	ctx context.Context
}

func NewEventEmitter(ctx context.Context) *EventEmitter {
	return &EventEmitter{
		ctx: ctx,
	}
}

func (e *EventEmitter) EmitEvent(event Event, payload interface{}) {
	runtime.EventsEmit(e.ctx, string(event), payload)
}
