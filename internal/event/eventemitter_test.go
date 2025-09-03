package event_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"gomander/internal/event"
	"gomander/internal/testutils/mocks"
)

func TestDefaultEventEmitter_EmitEvent(t *testing.T) {
	t.Run("Should emit event with payload", func(t *testing.T) {
		// Arrange
		eventKey := "test.event"
		eventPayload := "test payload"
		ctx := context.Background()
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		ee := event.NewDefaultEventEmitter(ctx, mockRuntimeFacade)

		mockRuntimeFacade.On("EventsEmit", ctx, eventKey, eventPayload).Return()

		// Act
		ee.EmitEvent(event.Event(eventKey), eventPayload)

		// Assert
		mock.AssertExpectationsForObjects(t, mockRuntimeFacade)
	})
}
