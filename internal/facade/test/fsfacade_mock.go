package test

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type MockFsFacade struct {
	mock.Mock
}

func (m *MockFsFacade) WriteFile(path string, data []byte, perm os.FileMode) error {
	args := m.Called(path, data, perm)
	return args.Error(0)
}

func (m *MockFsFacade) ReadFile(path string) ([]byte, error) {
	args := m.Called(path)
	return args.Get(0).([]byte), args.Error(1)
}
