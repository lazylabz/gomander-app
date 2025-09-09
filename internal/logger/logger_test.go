package logger_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"gomander/internal/facade/test"
	"gomander/internal/logger"
)

func TestDefaultLogger_Info(t *testing.T) {
	t.Run("Should call LogInfo with message", func(t *testing.T) {
		// Arrange
		message := "test info message"
		ctx := context.Background()
		mockRuntimeFacade := new(test.MockRuntimeFacade)

		l := logger.NewDefaultLogger(ctx, mockRuntimeFacade)

		mockRuntimeFacade.On("LogInfo", ctx, message).Return()

		// Act
		l.Info(message)

		// Assert
		mock.AssertExpectationsForObjects(t, mockRuntimeFacade)
	})
}

func TestDefaultLogger_Debug(t *testing.T) {
	t.Run("Should call LogDebug with message", func(t *testing.T) {
		// Arrange
		message := "test debug message"
		ctx := context.Background()
		mockRuntimeFacade := new(test.MockRuntimeFacade)

		l := logger.NewDefaultLogger(ctx, mockRuntimeFacade)

		mockRuntimeFacade.On("LogDebug", ctx, message).Return()

		// Act
		l.Debug(message)

		// Assert
		mock.AssertExpectationsForObjects(t, mockRuntimeFacade)
	})
}

func TestDefaultLogger_Error(t *testing.T) {
	t.Run("Should call LogError with message", func(t *testing.T) {
		// Arrange
		message := "test error message"
		ctx := context.Background()
		mockRuntimeFacade := new(test.MockRuntimeFacade)

		l := logger.NewDefaultLogger(ctx, mockRuntimeFacade)

		mockRuntimeFacade.On("LogError", ctx, message).Return()

		// Act
		l.Error(message)

		// Assert
		mock.AssertExpectationsForObjects(t, mockRuntimeFacade)
	})
}
