package test

import (
	"github.com/stretchr/testify/mock"

	"gomander/internal/event"
)

type MockEventEmitter struct {
	mock.Mock
}

func (m *MockEventEmitter) EmitEvent(event event.Event, payload interface{}) {
	m.Called(event, payload)
}
