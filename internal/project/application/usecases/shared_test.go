package usecases_test

import (
	"github.com/stretchr/testify/mock"

	"gomander/internal/config/domain"
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
