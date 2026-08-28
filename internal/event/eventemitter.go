package event

import (
	"context"
)

// EventSink is what carries an Event to whoever is listening on the other side
// of the desktop shell.
type EventSink interface {
	EventsEmit(ctx context.Context, eventName string, payload interface{})
}

type DefaultEventEmitter struct {
	ctx  context.Context
	sink EventSink
}

func NewDefaultEventEmitter(ctx context.Context, sink EventSink) *DefaultEventEmitter {
	return &DefaultEventEmitter{
		ctx:  ctx,
		sink: sink,
	}
}

func (e *DefaultEventEmitter) EmitEvent(event Event, payload interface{}) {
	e.sink.EventsEmit(e.ctx, string(event), payload)
}
