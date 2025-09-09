package handlers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	commanddomainevent "gomander/internal/command/domain/event"
	"gomander/internal/command/domain/test"
	"gomander/internal/commandgroup/application/handlers"
	test2 "gomander/internal/commandgroup/domain/test"
)

func TestDefaultAddCommandToGroupOnCommandDuplicated_GetEvent(t *testing.T) {
	// Arrange
	handler := handlers.NewAddCommandToGroupOnCommandDuplicated(nil, nil)

	// Act
	event := handler.GetEvent()

	// Assert
	_, ok := event.(commanddomainevent.CommandDuplicatedEvent)
	assert.True(t, ok)
}

func TestDefaultAddCommandToGroupOnCommandDuplicated(t *testing.T) {
	t.Run("Should do nothing if command was not duplicated inside a group", func(t *testing.T) {
		// Arrange
		mockCommandRepo := new(test.MockCommandRepository)
		mockCommandGroupRepo := new(test2.MockCommandGroupRepository)
		handler := handlers.NewAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)
		event := commanddomainevent.CommandDuplicatedEvent{
			CommandId:     "cmd-1",
			InsideGroupId: "",
		}

		// Act
		err := handler.Execute(event)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandRepo, mockCommandGroupRepo)
	})

	t.Run("Should add duplicated command to the group", func(t *testing.T) {
		// Arrange
		mockCommandRepo := new(test.MockCommandRepository)
		mockCommandGroupRepo := new(test2.MockCommandGroupRepository)
		handler := handlers.NewAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)
		event := commanddomainevent.CommandDuplicatedEvent{
			CommandId:     "cmd-1",
			InsideGroupId: "group-1",
		}

		existingGroup := test2.NewCommandGroupBuilder().Build()
		duplicatedCommand := test.NewCommandBuilder().WithId("cmd-1").Build()

		expectedUpdatedGroup := existingGroup
		expectedUpdatedGroup.Commands = append(expectedUpdatedGroup.Commands, duplicatedCommand)

		mockCommandGroupRepo.On("Get", event.InsideGroupId).Return(&existingGroup, nil)
		mockCommandRepo.On("Get", event.CommandId).Return(&duplicatedCommand, nil)
		mockCommandGroupRepo.On("Update", &expectedUpdatedGroup).Return(nil)

		// Act
		err := handler.Execute(event)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandRepo, mockCommandGroupRepo)
	})

	t.Run("Should do nothing if command is already in the group", func(t *testing.T) {
		// Arrange
		mockCommandRepo := new(test.MockCommandRepository)
		mockCommandGroupRepo := new(test2.MockCommandGroupRepository)
		handler := handlers.NewAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)
		event := commanddomainevent.CommandDuplicatedEvent{
			CommandId:     "cmd-1",
			InsideGroupId: "group-1",
		}

		duplicatedCommand := test.NewCommandBuilder().WithId("cmd-1").Build()
		existingGroup := test2.NewCommandGroupBuilder().WithCommands(duplicatedCommand).Build()

		mockCommandGroupRepo.On("Get", event.InsideGroupId).Return(&existingGroup, nil)

		// Act
		err := handler.Execute(event)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandRepo, mockCommandGroupRepo)
		mockCommandGroupRepo.AssertNotCalled(t, "Update", mock.Anything)
	})

	t.Run("Should do nothing if command is the wrong type", func(t *testing.T) {
		// Arrange
		mockCommandRepo := new(test.MockCommandRepository)
		mockCommandGroupRepo := new(test2.MockCommandGroupRepository)
		handler := handlers.NewAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)

		// Act
		err := handler.Execute(FakeEvent{})

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockCommandRepo, mockCommandGroupRepo)
	})

	t.Run("Should return error if failing to get command group", func(t *testing.T) {
		// Arrange
		mockCommandRepo := new(test.MockCommandRepository)
		mockCommandGroupRepo := new(test2.MockCommandGroupRepository)
		handler := handlers.NewAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)
		event := commanddomainevent.CommandDuplicatedEvent{
			CommandId:     "cmd-1",
			InsideGroupId: "group-1",
		}

		expectedError := errors.New("group not found")
		mockCommandGroupRepo.On("Get", event.InsideGroupId).Return(nil, expectedError)

		// Act
		err := handler.Execute(event)

		// Assert
		assert.ErrorIs(t, err, expectedError)
		mockCommandGroupRepo.AssertExpectations(t)
	})

	t.Run("Should return error if failing to get duplicated command", func(t *testing.T) {
		// Arrange
		mockCommandRepo := new(test.MockCommandRepository)
		mockCommandGroupRepo := new(test2.MockCommandGroupRepository)
		handler := handlers.NewAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)
		event := commanddomainevent.CommandDuplicatedEvent{
			CommandId:     "cmd-1",
			InsideGroupId: "group-1",
		}

		existingGroup := test2.NewCommandGroupBuilder().Build()
		mockCommandGroupRepo.On("Get", event.InsideGroupId).Return(&existingGroup, nil)

		expectedError := errors.New("command not found")
		mockCommandRepo.On("Get", event.CommandId).Return(nil, expectedError)

		// Act
		err := handler.Execute(event)

		// Assert
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t, mockCommandRepo, mockCommandGroupRepo)
	})

	t.Run("Should return error if failing to update command group", func(t *testing.T) {
		// Arrange
		mockCommandRepo := new(test.MockCommandRepository)
		mockCommandGroupRepo := new(test2.MockCommandGroupRepository)
		handler := handlers.NewAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)
		event := commanddomainevent.CommandDuplicatedEvent{
			CommandId:     "cmd-1",
			InsideGroupId: "group-1",
		}

		existingGroup := test2.NewCommandGroupBuilder().Build()
		duplicatedCommand := test.NewCommandBuilder().WithId("cmd-1").Build()

		expectedUpdatedGroup := existingGroup
		expectedUpdatedGroup.Commands = append(expectedUpdatedGroup.Commands, duplicatedCommand)

		expectedError := errors.New("update error")

		mockCommandGroupRepo.On("Get", event.InsideGroupId).Return(&existingGroup, nil)
		mockCommandRepo.On("Get", event.CommandId).Return(&duplicatedCommand, nil)
		mockCommandGroupRepo.On("Update", &expectedUpdatedGroup).Return(expectedError)

		// Act
		err := handler.Execute(event)

		// Assert
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t, mockCommandRepo, mockCommandGroupRepo)
	})
}
