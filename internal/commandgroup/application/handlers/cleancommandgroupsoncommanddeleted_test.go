package handlers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	commanddomainevent "gomander/internal/command/domain/event"
	"gomander/internal/commandgroup/application/handlers"
	commandgroupdomain "gomander/internal/commandgroup/domain"
)

type MockCommandGroupRepository struct {
	mock.Mock
}

func (m *MockCommandGroupRepository) Get(id string) (*commandgroupdomain.CommandGroup, error) {
	args := m.Called(id)
	return args.Get(0).(*commandgroupdomain.CommandGroup), args.Error(1)
}

func (m *MockCommandGroupRepository) GetAll(projectId string) ([]commandgroupdomain.CommandGroup, error) {
	args := m.Called(projectId)
	return args.Get(0).([]commandgroupdomain.CommandGroup), args.Error(1)
}

func (m *MockCommandGroupRepository) Create(commandGroup *commandgroupdomain.CommandGroup) error {
	args := m.Called(commandGroup)
	return args.Error(0)
}

func (m *MockCommandGroupRepository) Update(commandGroup *commandgroupdomain.CommandGroup) error {
	args := m.Called(commandGroup)
	return args.Error(0)
}

func (m *MockCommandGroupRepository) Delete(commandGroupId string) error {
	args := m.Called(commandGroupId)
	return args.Error(0)
}

func (m *MockCommandGroupRepository) RemoveCommandFromCommandGroups(commandId string) error {
	args := m.Called(commandId)
	return args.Error(0)
}

func (m *MockCommandGroupRepository) DeleteEmpty() error {
	args := m.Called()
	return args.Error(0)
}

type FakeEvent struct{}

func (FakeEvent) GetName() string { return "fake" }

var cmdId = "cmd-123"

func TestDefaultCleanCommandGroupsOnCommandDeleted(t *testing.T) {
	t.Run("Should remove command from command groups and delete empty groups", func(t *testing.T) {
		mockRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultCleanCommandGroupsOnCommandDeleted(mockRepo)
		event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

		mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(nil).Once()
		mockRepo.On("DeleteEmpty").Return(nil).Once()

		err := handler.Execute(event)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Should return error if failing to remove command from command groups", func(t *testing.T) {
		mockRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultCleanCommandGroupsOnCommandDeleted(mockRepo)
		event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

		expectedErr := errors.New("remove error")
		mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(expectedErr).Once()

		err := handler.Execute(event)
		assert.ErrorIs(t, err, expectedErr)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Should return error if failing to remove empty groups", func(t *testing.T) {
		mockRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultCleanCommandGroupsOnCommandDeleted(mockRepo)
		event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

		mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(nil).Once()
		expectedErr := errors.New("delete empty error")
		mockRepo.On("DeleteEmpty").Return(expectedErr).Once()

		err := handler.Execute(event)
		assert.ErrorIs(t, err, expectedErr)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Should do nothing if command is the wrong type", func(t *testing.T) {
		mockRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultCleanCommandGroupsOnCommandDeleted(mockRepo)

		err := handler.Execute(FakeEvent{})
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDefaultCleanCommandGroupsOnCommandDeleted_GetEvent(t *testing.T) {
	handler := handlers.NewDefaultCleanCommandGroupsOnCommandDeleted(nil)
	event := handler.GetEvent()
	_, ok := event.(commanddomainevent.CommandDeletedEvent)
	assert.True(t, ok)
}
