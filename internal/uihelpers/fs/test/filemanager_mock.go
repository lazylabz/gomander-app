package test

import (
	"github.com/stretchr/testify/mock"
)

type MockFileManager struct {
	mock.Mock
}

func (m *MockFileManager) OpenFolder(path string) error {
	args := m.Called(path)
	return args.Error(0)
}
