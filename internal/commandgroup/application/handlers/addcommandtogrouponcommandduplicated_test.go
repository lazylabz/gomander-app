package handlers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	commanddomain "gomander/internal/command/domain"
	commanddomainevent "gomander/internal/command/domain/event"
	"gomander/internal/commandgroup/application/handlers"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
	"gomander/internal/testutils"
)

func TestDefaultAddCommandToGroupOnCommandDuplicated_GetEvent(t *testing.T) {
	handler := handlers.NewDefaultAddCommandToGroupOnCommandDuplicated(nil, nil)
	event := handler.GetEvent()
	_, ok := event.(commanddomainevent.CommandDuplicatedEvent)
	assert.True(t, ok)
}

func TestDefaultAddCommandToGroupOnCommandDuplicated(t *testing.T) {
	t.Run("Should do nothing if command was not duplicated inside a group", func(t *testing.T) {
		mockCommandRepo := new(MockCommandRepository)
		mockCommandGroupRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)
		event := commanddomainevent.CommandDuplicatedEvent{
			CommandId:     "cmd-1",
			InsideGroupId: "",
		}

		err := handler.Execute(event)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandRepo, mockCommandGroupRepo)
	})
	t.Run("Should add duplicated command to the group", func(t *testing.T) {
		mockCommandRepo := new(MockCommandRepository)
		mockCommandGroupRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)
		event := commanddomainevent.CommandDuplicatedEvent{
			CommandId:     "cmd-1",
			InsideGroupId: "group-1",
		}

		existingGroupData := testutils.NewCommandGroup().Data()
		existingGroup := commandGroupDataToDomain(existingGroupData)
		duplicatedCommandData := testutils.NewCommand().WithId("cmd-1").Data()
		duplicatedCommand := commandDataToDomain(duplicatedCommandData)

		expectedUpdatedGroup := existingGroup
		expectedUpdatedGroup.Commands = append(expectedUpdatedGroup.Commands, duplicatedCommand)

		mockCommandGroupRepo.On("Get", event.InsideGroupId).Return(&existingGroup, nil)
		mockCommandRepo.On("Get", event.CommandId).Return(&duplicatedCommand, nil)
		mockCommandGroupRepo.On("Update", &expectedUpdatedGroup).Return(nil)

		err := handler.Execute(event)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandRepo, mockCommandGroupRepo)
	})
	t.Run("Should do nothing if command is already in the group", func(t *testing.T) {
		mockCommandRepo := new(MockCommandRepository)
		mockCommandGroupRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)
		event := commanddomainevent.CommandDuplicatedEvent{
			CommandId:     "cmd-1",
			InsideGroupId: "group-1",
		}

		duplicatedCommandData := testutils.NewCommand().WithId("cmd-1").Data()
		existingGroupData := testutils.NewCommandGroup().WithCommands(duplicatedCommandData).Data()
		existingGroup := commandGroupDataToDomain(existingGroupData)

		mockCommandGroupRepo.On("Get", event.InsideGroupId).Return(&existingGroup, nil)

		err := handler.Execute(event)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandRepo, mockCommandGroupRepo)
		mockCommandGroupRepo.AssertNotCalled(t, "Update", mock.Anything)
	})
	t.Run("Should do nothing if command is the wrong type", func(t *testing.T) {
		mockCommandRepo := new(MockCommandRepository)
		mockCommandGroupRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)

		err := handler.Execute(FakeEvent{})
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockCommandRepo, mockCommandGroupRepo)
	})
	t.Run("Should return error if failing to get command group", func(t *testing.T) {
		mockCommandRepo := new(MockCommandRepository)
		mockCommandGroupRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)
		event := commanddomainevent.CommandDuplicatedEvent{
			CommandId:     "cmd-1",
			InsideGroupId: "group-1",
		}

		expectedError := errors.New("group not found")
		mockCommandGroupRepo.On("Get", event.InsideGroupId).Return(nil, expectedError)

		err := handler.Execute(event)
		assert.ErrorIs(t, err, expectedError)
		mockCommandGroupRepo.AssertExpectations(t)
	})
	t.Run("Should return error if failing to get duplicated command", func(t *testing.T) {
		mockCommandRepo := new(MockCommandRepository)
		mockCommandGroupRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)
		event := commanddomainevent.CommandDuplicatedEvent{
			CommandId:     "cmd-1",
			InsideGroupId: "group-1",
		}

		existingGroupData := testutils.NewCommandGroup().Data()
		existingGroup := commandGroupDataToDomain(existingGroupData)
		mockCommandGroupRepo.On("Get", event.InsideGroupId).Return(&existingGroup, nil)

		expectedError := errors.New("command not found")
		mockCommandRepo.On("Get", event.CommandId).Return(nil, expectedError)

		err := handler.Execute(event)
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t, mockCommandRepo, mockCommandGroupRepo)
	})
	t.Run("Should return error if failing to update command group", func(t *testing.T) {
		mockCommandRepo := new(MockCommandRepository)
		mockCommandGroupRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultAddCommandToGroupOnCommandDuplicated(mockCommandRepo, mockCommandGroupRepo)
		event := commanddomainevent.CommandDuplicatedEvent{
			CommandId:     "cmd-1",
			InsideGroupId: "group-1",
		}

		existingGroupData := testutils.NewCommandGroup().Data()
		existingGroup := commandGroupDataToDomain(existingGroupData)
		duplicatedCommandData := testutils.NewCommand().WithId("cmd-1").Data()
		duplicatedCommand := commandDataToDomain(duplicatedCommandData)

		expectedUpdatedGroup := existingGroup
		expectedUpdatedGroup.Commands = append(expectedUpdatedGroup.Commands, duplicatedCommand)

		expectedError := errors.New("update error")

		mockCommandGroupRepo.On("Get", event.InsideGroupId).Return(&existingGroup, nil)
		mockCommandRepo.On("Get", event.CommandId).Return(&duplicatedCommand, nil)
		mockCommandGroupRepo.On("Update", &expectedUpdatedGroup).Return(expectedError)

		err := handler.Execute(event)
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t, mockCommandRepo, mockCommandGroupRepo)
	})
}

func commandDataToDomain(data testutils.CommandData) commanddomain.Command {
	return commanddomain.Command{
		Id:               data.Id,
		ProjectId:        data.ProjectId,
		Name:             data.Name,
		Command:          data.Command,
		WorkingDirectory: data.WorkingDirectory,
		Position:         data.Position,
	}
}

func commandGroupDataToDomain(data testutils.CommandGroupData) commandgroupdomain.CommandGroup {
	return commandgroupdomain.CommandGroup{
		Id:        data.Id,
		ProjectId: data.ProjectId,
		Name:      data.Name,
		Position:  data.Position,
		Commands:  array.Map(data.Commands, commandDataToDomain),
	}
}
