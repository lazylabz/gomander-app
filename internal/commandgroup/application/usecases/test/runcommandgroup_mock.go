package test

import "github.com/stretchr/testify/mock"

type MockRunCommandGroup struct {
	mock.Mock
}

func (m *MockRunCommandGroup) Execute(commandGroupId string) error {
	args := m.Called(commandGroupId)
	return args.Error(0)
}
