package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

func (m *MockRuntimeFacade) SaveFileDialog(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
	args := m.Called(ctx, dialogOptions)
	return args.String(0), args.Error(1)
}

func (m *MockRuntimeFacade) OpenFileDialog(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error) {
	args := m.Called(ctx, dialogOptions)
	return args.String(0), args.Error(1)
}
