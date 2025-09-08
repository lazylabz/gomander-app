package usecases_test

import (
	"github.com/stretchr/testify/mock"

	commanddomain "gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/config/domain"
	"gomander/internal/helpers/array"
	"gomander/internal/testutils"
)

type MockCommandRepository struct {
	mock.Mock
}

func (m *MockCommandRepository) Get(commandId string) (*commanddomain.Command, error) {
	args := m.Called(commandId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commanddomain.Command), args.Error(1)
}

func (m *MockCommandRepository) GetAll(projectId string) ([]commanddomain.Command, error) {
	args := m.Called(projectId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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

func (m *MockCommandRepository) Delete(commandId string) error {
	args := m.Called(commandId)
	return args.Error(0)
}

func (m *MockCommandRepository) DeleteAll(projectId string) error {
	args := m.Called(projectId)
	return args.Error(0)
}

type MockConfigRepository struct {
	mock.Mock
}

func (m *MockConfigRepository) GetOrCreate() (*domain.Config, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Config), args.Error(1)
}

func (m *MockConfigRepository) Update(config *domain.Config) error {
	args := m.Called(config)
	return args.Error(0)
}

func commandDataToDomain(data testutils.CommandData) commanddomain.Command {
	return commanddomain.Command{
		Id:               data.Id,
		ProjectId:        data.ProjectId,
		Name:             data.Name,
		Command:          data.Command,
		WorkingDirectory: data.WorkingDirectory,
		Position:         data.Position,
	}
}

func commandGroupDataToDomain(data testutils.CommandGroupData) commandgroupdomain.CommandGroup {
	return commandgroupdomain.CommandGroup{
		Id:        data.Id,
		ProjectId: data.ProjectId,
		Name:      data.Name,
		Position:  data.Position,
		Commands:  array.Map(data.Commands, commandDataToDomain),
	}
}
