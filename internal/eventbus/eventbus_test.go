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
	// Arrange
	bus := eventbus.NewInMemoryEventBus()
	handler := new(TestHandler)

	testEvent := &TestEvent{name: "test"}

	handler.On("GetEvent").Return(testEvent)
	handler.On("Execute", testEvent).Return(nil).Once()

	bus.RegisterHandler(handler)

	// Act
	evt := testEvent
	errs := bus.PublishSync(evt)

	// Assert
	assert.Empty(t, errs)
	mock.AssertExpectationsForObjects(t, handler)
}

func TestPublishSync_MultipleHandlers(t *testing.T) {
	// Arrange
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

	// Act
	errs := bus.PublishSync(evt)

	// Assert
	assert.Empty(t, errs)
	mock.AssertExpectationsForObjects(t, h1, h2)
}

func TestPublishSync_HandlerReturnsError(t *testing.T) {
	// Arrange
	bus := eventbus.NewInMemoryEventBus()
	evt := &TestEvent{name: "err"}

	handlerErr := errors.New("handler error")

	handler := new(TestHandler)
	handler.On("GetEvent").Return(evt)
	handler.On("Execute", evt).Return(handlerErr).Once()

	bus.RegisterHandler(handler)

	// Act
	errs := bus.PublishSync(evt)

	// Assert
	assert.Len(t, errs, 1)
	assert.Equal(t, handlerErr, errs[0])
}

func TestPublishSync_NoHandlers(t *testing.T) {
	// Arrange
	bus := eventbus.NewInMemoryEventBus()
	evt := &TestEvent{name: "nohandlers"}

	// Act
	errs := bus.PublishSync(evt)

	// Assert
	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
}

func TestCombined(t *testing.T) {
	t.Run("Should answer no error when no handler failed", func(t *testing.T) {
		// Act & Assert
		assert.NoError(t, eventbus.Combined("Errors occurred while doing something:", make([]error, 0)))
	})

	t.Run("Should list the failures under the summary", func(t *testing.T) {
		// Arrange
		errs := []error{errors.New("first went wrong"), errors.New("second went wrong")}

		// Act
		err := eventbus.Combined("Errors occurred while doing something:", errs)

		// Assert
		assert.EqualError(
			t,
			err,
			"Errors occurred while doing something:\n- first went wrong\n- second went wrong",
		)
	})
}
