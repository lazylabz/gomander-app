package eventbus

type Event interface {
	GetName() string
}

type EventHandler interface {
	Execute(Event) error
	GetEvent() Event
}

type EventBus interface {
	RegisterHandler(EventHandler)
	PublishSync(Event) []error
}

type InMemoryEventBus struct {
	eventHandlers map[string][]EventHandler
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		eventHandlers: make(map[string][]EventHandler),
	}
}

func (e *InMemoryEventBus) RegisterHandler(handler EventHandler) {
	eventName := handler.GetEvent().GetName()
	e.eventHandlers[eventName] = append(e.eventHandlers[eventName], handler)
}

// Combined answers the handler failures PublishSync reported as one error,
// listed under summary, or nil when there were none. Every operation that
// publishes an event has to report what its handlers did to the user, and this
// is the one place that decides how that reads.
//
// The handler errors stay in the chain: a caller reads the summary, but has to
// be able to tell a missing entity from a storage failure with errors.Is.
func Combined(summary string, errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	message := summary

	for _, err := range errs {
		message += "\n- " + err.Error()
	}

	// Copied and returned by pointer: errors.New answered a comparable pointer,
	// and a value holding a slice would panic on an == between two of these.
	return &combinedError{message: message, errs: append([]error(nil), errs...)}
}

func (e *InMemoryEventBus) PublishSync(event Event) []error {
	errs := make([]error, 0)
	if handlers, ok := e.eventHandlers[event.GetName()]; ok {
		ch := make(chan error, len(handlers))
		for _, handler := range handlers {
			go func(h EventHandler) {
				ch <- h.Execute(event)
			}(handler)
		}

		for range handlers {
			if err := <-ch; err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

type combinedError struct {
	message string
	errs    []error
}

func (e *combinedError) Error() string { return e.message }

func (e *combinedError) Unwrap() []error { return e.errs }

// Not used for now
//
//	func (e *InMemoryEventBus) PublishAsync(event Event) {
//		if handlers, ok := e.eventHandlers[event.GetName()]; ok {
//			for _, handler := range handlers {
//				go func(h func(event Event) error) {
//					if err := h(event); err != nil {
//						// Log the error
//						log.Printf("Error handling event %s: %v", event.GetName(), err)
//					}
//				}(handler)
//			}
//		}
//	}
