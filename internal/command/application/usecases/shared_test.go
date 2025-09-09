package usecases_test

import (
	"github.com/stretchr/testify/mock"

	"gomander/internal/eventbus"
)

type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) RegisterHandler(handler eventbus.EventHandler) {
	m.Called(handler)
}

func (m *MockEventBus) PublishSync(e eventbus.Event) []error {
	args := m.Called(e)
	return args.Get(0).([]error)
}
