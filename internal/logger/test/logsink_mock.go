package test

import (
	"github.com/stretchr/testify/mock"
)

type MockLogSink struct {
	mock.Mock
}

func (m *MockLogSink) LogInfo(message string) {
	m.Called(message)
}

func (m *MockLogSink) LogDebug(message string) {
	m.Called(message)
}

func (m *MockLogSink) LogError(message string) {
	m.Called(message)
}
