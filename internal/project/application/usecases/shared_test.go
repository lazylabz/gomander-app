package usecases_test

import (
	"os"

	"github.com/stretchr/testify/mock"

	commanddomain "gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/config/domain"
	"gomander/internal/eventbus"
	projectdomain "gomander/internal/project/domain"
)

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

type MockCommandGroupRepository struct {
	mock.Mock
}

func (m *MockCommandGroupRepository) Get(id string) (*commandgroupdomain.CommandGroup, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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

func (m *MockCommandGroupRepository) DeleteEmpty() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCommandGroupRepository) DeleteAll(projectId string) error {
	args := m.Called(projectId)
	return args.Error(0)
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
