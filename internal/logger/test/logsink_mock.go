package test

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockLogSink struct {
	mock.Mock
}

func (m *MockLogSink) LogInfo(ctx context.Context, message string) {
	m.Called(ctx, message)
}

func (m *MockLogSink) LogDebug(ctx context.Context, message string) {
	m.Called(ctx, message)
}

func (m *MockLogSink) LogError(ctx context.Context, message string) {
	m.Called(ctx, message)
}
