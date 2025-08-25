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

var cmdId = "cmd-123"

func TestDefaultCleanCommandGroupsOnCommandDeleted_Success(t *testing.T) {
	mockRepo := new(MockCommandGroupRepository)
	handler := handlers.NewDefaultCleanCommandGroupsOnCommandDeleted(mockRepo)
	event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

	mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(nil).Once()
	mockRepo.On("DeleteEmpty").Return(nil).Once()

	err := handler.Execute(event)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDefaultCleanCommandGroupsOnCommandDeleted_RemoveCommandFails(t *testing.T) {
	mockRepo := new(MockCommandGroupRepository)
	handler := handlers.NewDefaultCleanCommandGroupsOnCommandDeleted(mockRepo)
	event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

	expectedErr := errors.New("remove error")
	mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(expectedErr).Once()

	err := handler.Execute(event)
	assert.ErrorIs(t, err, expectedErr)
	mockRepo.AssertExpectations(t)
}

func TestDefaultCleanCommandGroupsOnCommandDeleted_DeleteEmptyFails(t *testing.T) {
	mockRepo := new(MockCommandGroupRepository)
	handler := handlers.NewDefaultCleanCommandGroupsOnCommandDeleted(mockRepo)
	event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

	mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(nil).Once()
	expectedErr := errors.New("delete empty error")
	mockRepo.On("DeleteEmpty").Return(expectedErr).Once()

	err := handler.Execute(event)
	assert.ErrorIs(t, err, expectedErr)
	mockRepo.AssertExpectations(t)
}

type fakeEvent struct{}

func (fakeEvent) GetName() string { return "fake" }

func TestDefaultCleanCommandGroupsOnCommandDeleted_InvalidEventType(t *testing.T) {
	mockRepo := new(MockCommandGroupRepository)
	handler := handlers.NewDefaultCleanCommandGroupsOnCommandDeleted(mockRepo)

	err := handler.Execute(fakeEvent{})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
