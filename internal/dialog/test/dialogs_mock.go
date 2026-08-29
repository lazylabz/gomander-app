package test

import (
	"github.com/stretchr/testify/mock"

	"gomander/internal/dialog"
)

type MockDialogs struct {
	mock.Mock
}

func (m *MockDialogs) AskForFileToOpen(request dialog.OpenFileRequest) (string, error) {
	args := m.Called(request)
	return args.String(0), args.Error(1)
}

func (m *MockDialogs) AskWhereToSaveFile(request dialog.SaveFileRequest) (string, error) {
	args := m.Called(request)
	return args.String(0), args.Error(1)
}

func (m *MockDialogs) AskForDirectory(request dialog.PickDirectoryRequest) (string, error) {
	args := m.Called(request)
	return args.String(0), args.Error(1)
}
