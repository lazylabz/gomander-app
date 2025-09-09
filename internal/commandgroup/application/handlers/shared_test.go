package handlers_test

import (
	"github.com/stretchr/testify/mock"

	commandgroupdomain "gomander/internal/commandgroup/domain"
)

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

type FakeEvent struct{}

func (FakeEvent) GetName() string { return "fake" }
