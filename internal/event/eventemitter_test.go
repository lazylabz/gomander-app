package event_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"gomander/internal/event"
	"gomander/internal/event/test"
)

func TestDefaultEventEmitter_EmitEvent(t *testing.T) {
	t.Run("Should emit event with payload", func(t *testing.T) {
		// Arrange
		eventKey := "test.event"
		eventPayload := "test payload"
		mockEventSink := new(test.MockEventSink)

		ee := event.NewDefaultEventEmitter(mockEventSink)

		mockEventSink.On("EventsEmit", eventKey, eventPayload).Return()

		// Act
		ee.EmitEvent(event.Event(eventKey), eventPayload)

		// Assert
		mock.AssertExpectationsForObjects(t, mockEventSink)
	})
}
