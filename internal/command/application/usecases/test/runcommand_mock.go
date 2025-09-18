package test

import "github.com/stretchr/testify/mock"

type MockRunCommands struct {
	mock.Mock
}

func (m *MockRunCommands) Execute(commandId string) error {
	args := m.Called(commandId)
	return args.Error(0)
}
