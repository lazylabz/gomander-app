package test

import (
	"sync"

	"gomander/internal/eventbus"
)

type DomainEventsRecorder struct {
	mutex  sync.Mutex
	events []eventbus.Event
}

func (r *DomainEventsRecorder) PublishSync(e eventbus.Event) []error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.events = append(r.events, e)

	return nil
}

func (r *DomainEventsRecorder) Published() []eventbus.Event {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return append([]eventbus.Event(nil), r.events...)
}
