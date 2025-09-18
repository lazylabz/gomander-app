package test

import "github.com/stretchr/testify/mock"

type MockGetRunningCommandIds struct {
	mock.Mock
}

func (m *MockGetRunningCommandIds) Execute() []string {
	args := m.Called()
	return args.Get(0).([]string)
}
