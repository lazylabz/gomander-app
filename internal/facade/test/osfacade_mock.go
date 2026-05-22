package test

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type MockOSFacade struct {
	mock.Mock
}

func (m *MockOSFacade) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(os.FileInfo), args.Error(1)
}

func (m *MockOSFacade) TempDir() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockOSFacade) Create(name string) (*os.File, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*os.File), args.Error(1)
}
