package app_test

import (
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(message string) {
	m.Called(message)
}

func (m *MockLogger) Debug(message string) {
	m.Called(message)
}

func (m *MockLogger) Error(message string) {
	m.Called(message)
}
