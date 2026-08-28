package test

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRuntimeFacade struct {
	mock.Mock
}

func (m *MockRuntimeFacade) EventsEmit(ctx context.Context, eventName string, optionalData interface{}) {
	m.Called(ctx, eventName, optionalData)
}

func (m *MockRuntimeFacade) LogInfo(ctx context.Context, message string) {
	m.Called(ctx, message)
}

func (m *MockRuntimeFacade) LogDebug(ctx context.Context, message string) {
	m.Called(ctx, message)
}

func (m *MockRuntimeFacade) LogError(ctx context.Context, message string) {
	m.Called(ctx, message)
}

func (m *MockRuntimeFacade) OpenFolderInFileManager(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockRuntimeFacade) CloseApp(ctx context.Context) {
	m.Called(ctx)
}
