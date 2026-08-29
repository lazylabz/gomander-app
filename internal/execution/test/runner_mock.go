package test

import (
	"github.com/stretchr/testify/mock"

	commanddomain "gomander/internal/command/domain"
	"gomander/internal/execution"
)

type MockRunner struct {
	mock.Mock
}

func (m *MockRunner) RunCommand(command *commanddomain.Command, environment execution.Environment) error {
	args := m.Called(command, environment)
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
