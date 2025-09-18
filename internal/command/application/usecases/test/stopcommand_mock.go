package test

import "github.com/stretchr/testify/mock"

type MockStopCommand struct {
	mock.Mock
}

func (m *MockStopCommand) Execute(commandId string) error {
	args := m.Called(commandId)
	return args.Error(0)
}
