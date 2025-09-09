package usecases_test

import (
	"github.com/stretchr/testify/mock"

	commanddomain "gomander/internal/command/domain"
	"gomander/internal/config/domain"
	"gomander/internal/eventbus"
	projectdomain "gomander/internal/project/domain"
)

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
