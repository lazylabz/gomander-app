package event

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type EventEmitter interface {
	EmitEvent(event Event, payload interface{})
}

type DefaultEventEmitter struct {
	ctx context.Context
}

func NewDefaultEventEmitter(ctx context.Context) *DefaultEventEmitter {
	return &DefaultEventEmitter{
		ctx: ctx,
	}
}

func (e *DefaultEventEmitter) EmitEvent(event Event, payload interface{}) {
	runtime.EventsEmit(e.ctx, string(event), payload)
}
