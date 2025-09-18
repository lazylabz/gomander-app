package test

import (
	"github.com/stretchr/testify/mock"

	"gomander/internal/commandgroup/domain"
)

type MockGetCommandGroups struct {
	mock.Mock
}

func (m *MockGetCommandGroups) Execute() ([]domain.CommandGroup, error) {
	args := m.Called()
	return args.Get(0).([]domain.CommandGroup), args.Error(1)
}
