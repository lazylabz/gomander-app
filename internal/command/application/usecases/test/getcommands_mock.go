package test

import (
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/domain"
)

type MockGetCommands struct {
	mock.Mock
}

func (g *MockGetCommands) Execute() ([]domain.Command, error) {
	args := g.Called()
	return args.Get(0).([]domain.Command), args.Error(1)
}
