package test

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type MockIOFacade struct {
	mock.Mock
}

func (m *MockIOFacade) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	args := m.Called(dst, src)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockIOFacade) ReadAll(r io.Reader) ([]byte, error) {
	args := m.Called(r)
	return args.Get(0).([]byte), args.Error(1)
}
