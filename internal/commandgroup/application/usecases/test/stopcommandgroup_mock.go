package test

import "github.com/stretchr/testify/mock"

type MockStopCommandGroup struct {
	mock.Mock
}

func (m *MockStopCommandGroup) Execute(commandGroupId string) error {
	args := m.Called(commandGroupId)
	return args.Error(0)
}
