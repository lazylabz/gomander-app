package logger_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"gomander/internal/logger"
	"gomander/internal/testutils/mocks"
)

func TestDefaultLogger_Info(t *testing.T) {
	t.Run("Should call LogInfo with message", func(t *testing.T) {
		message := "test info message"
		ctx := context.Background()
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		l := logger.NewDefaultLogger(ctx, mockRuntimeFacade)

		mockRuntimeFacade.On("LogInfo", ctx, message).Return()

		l.Info(message)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade)
	})
}

func TestDefaultLogger_Debug(t *testing.T) {
	t.Run("Should call LogDebug with message", func(t *testing.T) {
		message := "test debug message"
		ctx := context.Background()
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		l := logger.NewDefaultLogger(ctx, mockRuntimeFacade)

		mockRuntimeFacade.On("LogDebug", ctx, message).Return()

		l.Debug(message)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade)
	})
}

func TestDefaultLogger_Error(t *testing.T) {
	t.Run("Should call LogError with message", func(t *testing.T) {
		message := "test error message"
		ctx := context.Background()
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		l := logger.NewDefaultLogger(ctx, mockRuntimeFacade)

		mockRuntimeFacade.On("LogError", ctx, message).Return()

		l.Error(message)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade)
	})
}
