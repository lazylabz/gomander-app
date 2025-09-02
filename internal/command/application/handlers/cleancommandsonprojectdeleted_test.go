package handlers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/application/handlers"
	commanddomain "gomander/internal/command/domain"
	projectdomainevent "gomander/internal/project/domain/event"
)

type MockCommandRepository struct {
	mock.Mock
}

func (m *MockCommandRepository) Get(commandId string) (*commanddomain.Command, error) {
	args := m.Called(commandId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commanddomain.Command), args.Error(1)
}

func (m *MockCommandRepository) GetAll(projectId string) ([]commanddomain.Command, error) {
	args := m.Called(projectId)
	return args.Get(0).([]commanddomain.Command), args.Error(1)
}

func (m *MockCommandRepository) Create(command *commanddomain.Command) error {
	args := m.Called(command)
	return args.Error(0)
}

func (m *MockCommandRepository) Update(command *commanddomain.Command) error {
	args := m.Called(command)
	return args.Error(0)
}

func (m *MockCommandRepository) Delete(commandId string) error {
	args := m.Called(commandId)
	return args.Error(0)
}

func (m *MockCommandRepository) DeleteAll(projectId string) error {
	args := m.Called(projectId)
	return args.Error(0)
}

type FakeEvent struct{}

func (FakeEvent) GetName() string { return "fake" }

var pjId = "pj-123"

func TestDefaultCleanCommandsOnProjectDeleted(t *testing.T) {
	t.Run("Should remove command from command groups and delete empty groups", func(t *testing.T) {
		mockRepo := new(MockCommandRepository)
		handler := handlers.NewCleanCommandOnProjectDeleted(mockRepo)
		event := projectdomainevent.ProjectDeletedEvent{ProjectId: pjId}

		mockRepo.On("DeleteAll", pjId).Return(nil).Once()

		err := handler.Execute(event)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Should return error if failing to remove command from command groups", func(t *testing.T) {
		mockRepo := new(MockCommandRepository)
		handler := handlers.NewCleanCommandOnProjectDeleted(mockRepo)
		event := projectdomainevent.ProjectDeletedEvent{ProjectId: pjId}

		expectedErr := errors.New("remove error")
		mockRepo.On("DeleteAll", pjId).Return(expectedErr).Once()

		err := handler.Execute(event)
		assert.ErrorIs(t, err, expectedErr)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Should do nothing if command is the wrong type", func(t *testing.T) {
		mockRepo := new(MockCommandRepository)
		handler := handlers.NewCleanCommandOnProjectDeleted(mockRepo)

		err := handler.Execute(FakeEvent{})
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDefaultCleanCommandsOnProjectDeleted_GetEvent(t *testing.T) {
	handler := handlers.NewCleanCommandOnProjectDeleted(nil)
	event := handler.GetEvent()
	_, ok := event.(projectdomainevent.ProjectDeletedEvent)
	assert.True(t, ok)
}
