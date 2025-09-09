package test

import (
	"github.com/stretchr/testify/mock"

	commanddomain "gomander/internal/command/domain"
)

type MockRunner struct {
	mock.Mock
}

func (m *MockRunner) RunCommands(commands []commanddomain.Command, environmentPaths []string, baseWorkingDirectory string) error {
	args := m.Called(commands, environmentPaths, baseWorkingDirectory)
	return args.Error(0)
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
