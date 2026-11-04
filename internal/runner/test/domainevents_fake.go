package test

import (
	"sync"

	"gomander/internal/eventbus"
)

type DomainEventsRecorder struct {
	mutex  sync.Mutex
	events []eventbus.Event
	// HandlerFailure, when set, is what every handler of a published event reports
	HandlerFailure error
}

func (r *DomainEventsRecorder) PublishSync(e eventbus.Event) []error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.events = append(r.events, e)

	if r.HandlerFailure != nil {
		return []error{r.HandlerFailure}
	}

	return nil
}

func (r *DomainEventsRecorder) Published() []eventbus.Event {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return append([]eventbus.Event(nil), r.events...)
}
