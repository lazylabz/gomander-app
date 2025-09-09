package test

import (
	"github.com/stretchr/testify/mock"

	commanddomain "gomander/internal/command/domain"
)

type MockCommandRepository struct {
	mock.Mock
}

func (m *MockCommandRepository) Get(id string) (*commanddomain.Command, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commanddomain.Command), args.Error(1)
}

func (m *MockCommandRepository) GetAll(projectId string) ([]commanddomain.Command, error) {
	args := m.Called(projectId)
	return args.Get(0).([]commanddomain.Command), args.Error(1)
}

func (m *MockCommandRepository) Create(command *commanddomain.Command) error {
	args := m.Called(command)
	return args.Error(0)
}

func (m *MockCommandRepository) Update(command *commanddomain.Command) error {
	args := m.Called(command)
	return args.Error(0)
}

func (m *MockCommandRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCommandRepository) DeleteAll(projectId string) error {
	args := m.Called(projectId)
	return args.Error(0)
}
