package test

import (
	"github.com/stretchr/testify/mock"

	commanddomain "gomander/internal/command/domain"
)

type MockRunner struct {
	mock.Mock
}

func (m *MockRunner) StopRunningCommands(commands []commanddomain.Command) error {
	args := m.Called(commands)
	return args.Error(0)
}

func (m *MockRunner) RunCommands(commands []commanddomain.Command, environmentPaths []string, baseWorkingDirectory string, failurePatterns []string) error {
	args := m.Called(commands, environmentPaths, baseWorkingDirectory, failurePatterns)
	return args.Error(0)
}

func (m *MockRunner) RunCommand(command *commanddomain.Command, environmentPaths []string, baseWorkingDirectory string, failurePatterns []string) error {
	args := m.Called(command, environmentPaths, baseWorkingDirectory, failurePatterns)
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

func (m *MockRunner) GetRunningCommandIds() []string {
	args := m.Called()
	return args.Get(0).([]string)
}
