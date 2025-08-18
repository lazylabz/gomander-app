package event

import (
	"context"

	"gomander/internal/facade"
)

type EventEmitter interface {
	EmitEvent(event Event, payload interface{})
}

type DefaultEventEmitter struct {
	ctx     context.Context
	runtime facade.RuntimeFacade
}

func NewDefaultEventEmitter(ctx context.Context, runtime facade.RuntimeFacade) *DefaultEventEmitter {
	return &DefaultEventEmitter{
		ctx:     ctx,
		runtime: runtime,
	}
}

func (e *DefaultEventEmitter) EmitEvent(event Event, payload interface{}) {
	e.runtime.EventsEmit(e.ctx, string(event), payload)
}
