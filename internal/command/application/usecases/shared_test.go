package usecases_test

import (
	"github.com/stretchr/testify/mock"

	commanddomain "gomander/internal/command/domain"
	"gomander/internal/eventbus"
)

type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) RegisterHandler(handler eventbus.EventHandler) {
	m.Called(handler)
}

func (m *MockEventBus) PublishSync(e eventbus.Event) []error {
	args := m.Called(e)
	return args.Get(0).([]error)
}

type MockRunner struct {
	mock.Mock
}

func (m *MockRunner) RunCommand(command *commanddomain.Command, environmentPaths []string, baseWorkingDirectory string) error {
	args := m.Called(command, environmentPaths, baseWorkingDirectory)
	return args.Error(0)
}

func (m *MockRunner) StopRunningCommand(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRunner) StopAllRunningCommands() []error {
	args := m.Called()
	return args.Get(0).([]error)
}
