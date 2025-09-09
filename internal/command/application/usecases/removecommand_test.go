package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/application/usecases"
	commanddomainevent "gomander/internal/command/domain/event"
	"gomander/internal/command/domain/test"
)

func TestDefaultRemoveCommand_Execute(t *testing.T) {
	t.Run("Should remove the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockEventBus := new(MockEventBus)

		sut := usecases.NewRemoveCommand(mockCommandRepository, mockEventBus)

		commandId := "command1"

		mockCommandRepository.On("Delete", commandId).Return(nil)
		mockEventBus.On("PublishSync", commanddomainevent.NewCommandDeletedEvent(commandId)).Return(make([]error, 0))

		// Act
		err := sut.Execute(commandId)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockEventBus,
		)
	})

	t.Run("Should return an error if fails to remove the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockEventBus := new(MockEventBus)

		sut := usecases.NewRemoveCommand(mockCommandRepository, mockEventBus)

		commandId := "command1"

		mockCommandRepository.On("Delete", commandId).Return(errors.New("failed to delete command"))

		// Act
		err := sut.Execute(commandId)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
		)
	})

	t.Run("Should return an error if side effect fail", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockEventBus := new(MockEventBus)

		sut := usecases.NewRemoveCommand(mockCommandRepository, mockEventBus)

		commandId := "command1"

		mockCommandRepository.On("Delete", commandId).Return(nil)
		mockEventBus.On("PublishSync", commanddomainevent.NewCommandDeletedEvent(commandId)).Return([]error{errors.New("Something happened")})

		// Act
		err := sut.Execute(commandId)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockEventBus,
		)
	})
}
