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
