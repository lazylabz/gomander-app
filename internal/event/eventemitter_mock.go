package event

import "github.com/stretchr/testify/mock"

type MockEventEmitter struct {
	mock.Mock
}

func (m *MockEventEmitter) EmitEvent(event Event, payload interface{}) {
	m.Called(event, payload)
}
