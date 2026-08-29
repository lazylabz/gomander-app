package test

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockEventSink struct {
	mock.Mock
}

func (m *MockEventSink) EventsEmit(ctx context.Context, eventName string, payload interface{}) {
	m.Called(ctx, eventName, payload)
}
