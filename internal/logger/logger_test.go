package logger_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"gomander/internal/logger"
	"gomander/internal/logger/test"
)

func TestDefaultLogger_Info(t *testing.T) {
	t.Run("Should call LogInfo with message", func(t *testing.T) {
		// Arrange
		message := "test info message"
		ctx := context.Background()
		mockLogSink := new(test.MockLogSink)

		l := logger.NewDefaultLogger(ctx, mockLogSink)

		mockLogSink.On("LogInfo", ctx, message).Return()

		// Act
		l.Info(message)

		// Assert
		mock.AssertExpectationsForObjects(t, mockLogSink)
	})
}

func TestDefaultLogger_Debug(t *testing.T) {
	t.Run("Should call LogDebug with message", func(t *testing.T) {
		// Arrange
		message := "test debug message"
		ctx := context.Background()
		mockLogSink := new(test.MockLogSink)

		l := logger.NewDefaultLogger(ctx, mockLogSink)

		mockLogSink.On("LogDebug", ctx, message).Return()

		// Act
		l.Debug(message)

		// Assert
		mock.AssertExpectationsForObjects(t, mockLogSink)
	})
}

func TestDefaultLogger_Error(t *testing.T) {
	t.Run("Should call LogError with message", func(t *testing.T) {
		// Arrange
		message := "test error message"
		ctx := context.Background()
		mockLogSink := new(test.MockLogSink)

		l := logger.NewDefaultLogger(ctx, mockLogSink)

		mockLogSink.On("LogError", ctx, message).Return()

		// Act
		l.Error(message)

		// Assert
		mock.AssertExpectationsForObjects(t, mockLogSink)
	})
}
