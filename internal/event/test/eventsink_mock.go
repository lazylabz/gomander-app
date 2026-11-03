package test

import (
	"github.com/stretchr/testify/mock"
)

type MockEventSink struct {
	mock.Mock
}

func (m *MockEventSink) EventsEmit(eventName string, payload interface{}) {
	m.Called(eventName, payload)
}
