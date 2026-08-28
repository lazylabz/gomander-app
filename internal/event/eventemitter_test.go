package event_test

import (
	"context"
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
		ctx := context.Background()
		mockEventSink := new(test.MockEventSink)

		ee := event.NewDefaultEventEmitter(ctx, mockEventSink)

		mockEventSink.On("EventsEmit", ctx, eventKey, eventPayload).Return()

		// Act
		ee.EmitEvent(event.Event(eventKey), eventPayload)

		// Assert
		mock.AssertExpectationsForObjects(t, mockEventSink)
	})
}
