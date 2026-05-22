package test

import (
	"github.com/stretchr/testify/mock"
)

type MockOpenFacade struct {
	mock.Mock
}

func (m *MockOpenFacade) Run(input string) error {
	args := m.Called(input)
	return args.Error(0)
}
