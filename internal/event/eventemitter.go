package event

// EventSink is what carries an Event to whoever is listening on the other side
// of the desktop shell.
type EventSink interface {
	EventsEmit(eventName string, payload interface{})
}

type DefaultEventEmitter struct {
	sink EventSink
}

func NewDefaultEventEmitter(sink EventSink) *DefaultEventEmitter {
	return &DefaultEventEmitter{
		sink: sink,
	}
}

func (e *DefaultEventEmitter) EmitEvent(event Event, payload interface{}) {
	e.sink.EventsEmit(string(event), payload)
}
