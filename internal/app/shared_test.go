package app_test

import (
	"github.com/stretchr/testify/mock"

	commanddomain "gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/config/domain"
	"gomander/internal/event"
	"gomander/internal/eventbus"
	"gomander/internal/helpers/array"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/testutils"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(message string) {
	m.Called(message)
}

func (m *MockLogger) Debug(message string) {
	m.Called(message)
}

func (m *MockLogger) Error(message string) {
	m.Called(message)
}

type MockEventEmitter struct {
	mock.Mock
}

func (m *MockEventEmitter) EmitEvent(event event.Event, payload interface{}) {
	m.Called(event, payload)
}

type MockUserConfigRepository struct {
	mock.Mock
}

func (m *MockUserConfigRepository) GetOrCreate() (*domain.Config, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Config), args.Error(1)
}

func (m *MockUserConfigRepository) Update(config *domain.Config) error {
	args := m.Called(config)
	return args.Error(0)
}

type MockCommandGroupRepository struct {
	mock.Mock
}

func (m *MockCommandGroupRepository) Get(id string) (*commandgroupdomain.CommandGroup, error) {
	args := m.Called(id)
	return args.Get(0).(*commandgroupdomain.CommandGroup), args.Error(1)
}

func (m *MockCommandGroupRepository) GetAll(projectId string) ([]commandgroupdomain.CommandGroup, error) {
	args := m.Called(projectId)
	return args.Get(0).([]commandgroupdomain.CommandGroup), args.Error(1)
}

func (m *MockCommandGroupRepository) Create(commandGroup *commandgroupdomain.CommandGroup) error {
	args := m.Called(commandGroup)
	return args.Error(0)
}

func (m *MockCommandGroupRepository) Update(commandGroup *commandgroupdomain.CommandGroup) error {
	args := m.Called(commandGroup)
	return args.Error(0)
}

func (m *MockCommandGroupRepository) Delete(commandGroupId string) error {
	args := m.Called(commandGroupId)
	return args.Error(0)
}

func (m *MockCommandGroupRepository) RemoveCommandFromCommandGroups(commandId string) error {
	args := m.Called(commandId)
	return args.Error(0)
}

func (m *MockCommandGroupRepository) GetEmptyCommandGroups() ([]commandgroupdomain.CommandGroup, error) {
	args := m.Called()
	return args.Get(0).([]commandgroupdomain.CommandGroup), args.Error(1)
}

type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) GetAll() ([]projectdomain.Project, error) {
	args := m.Called()
	return args.Get(0).([]projectdomain.Project), args.Error(1)
}

func (m *MockProjectRepository) Get(id string) (*projectdomain.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectdomain.Project), args.Error(1)
}

func (m *MockProjectRepository) Create(project projectdomain.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) Update(project projectdomain.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockRunner struct {
	mock.Mock
}

func (m *MockRunner) RunCommand(command *commanddomain.Command, environmentPaths []string, baseWorkingDirectory string) error {
	args := m.Called(command, environmentPaths, baseWorkingDirectory)
	return args.Error(0)
}

func (m *MockRunner) StopRunningCommand(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRunner) StopAllRunningCommands() []error {
	args := m.Called()
	return args.Get(0).([]error)
}

type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) RegisterHandler(handler eventbus.EventHandler) {
	m.Called(handler)
}

func (m *MockEventBus) PublishSync(e eventbus.Event) []error {
	args := m.Called(e)
	return args.Get(0).([]error)
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
