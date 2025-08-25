package eventbus_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/eventbus"
)

type TestEvent struct {
	name string
	data string
}

func (e *TestEvent) GetName() string {
	return e.name
}

type TestHandler struct {
	mock.Mock
}

func (t *TestHandler) Execute(event eventbus.Event) error {
	args := t.Called(event)
	return args.Error(0)
}

func (t *TestHandler) GetEvent() eventbus.Event {
	args := t.Called()
	return args.Get(0).(eventbus.Event)
}

func TestRegisterHandlerAndPublishSync_Success(t *testing.T) {
	bus := eventbus.NewInMemoryEventBus()
	handler := new(TestHandler)

	testEvent := &TestEvent{name: "test"}

	handler.On("GetEvent").Return(testEvent)
	handler.On("Execute", testEvent).Return(nil).Once()

	bus.RegisterHandler(handler)
	evt := testEvent
	errs := bus.PublishSync(evt)
	assert.Empty(t, errs)
	mock.AssertExpectationsForObjects(t, handler)
}

func TestPublishSync_MultipleHandlers(t *testing.T) {
	bus := eventbus.NewInMemoryEventBus()
	evt := &TestEvent{name: "multi"}

	h1 := new(TestHandler)
	h2 := new(TestHandler)

	h1.On("GetEvent").Return(evt)
	h1.On("Execute", evt).Return(nil).Once()
	h2.On("GetEvent").Return(evt)
	h2.On("Execute", evt).Return(nil).Once()

	bus.RegisterHandler(h1)
	bus.RegisterHandler(h2)
	errs := bus.PublishSync(evt)
	assert.Empty(t, errs)

	mock.AssertExpectationsForObjects(t, h1, h2)
}

func TestPublishSync_HandlerReturnsError(t *testing.T) {
	bus := eventbus.NewInMemoryEventBus()
	evt := &TestEvent{name: "err"}

	handlerErr := errors.New("handler error")

	handler := new(TestHandler)
	handler.On("GetEvent").Return(evt)
	handler.On("Execute", evt).Return(handlerErr).Once()

	bus.RegisterHandler(handler)
	errs := bus.PublishSync(evt)
	assert.Len(t, errs, 1)
	assert.Equal(t, handlerErr, errs[0])
}

func TestPublishSync_NoHandlers(t *testing.T) {
	bus := eventbus.NewInMemoryEventBus()
	evt := &TestEvent{name: "nohandlers"}
	errs := bus.PublishSync(evt)
	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
}
