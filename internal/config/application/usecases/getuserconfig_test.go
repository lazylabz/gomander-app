package usecases_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/config/application/usecases"
	"gomander/internal/config/domain"
)

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

func TestApp_GetUserConfig(t *testing.T) {
	mockRepository := new(MockUserConfigRepository)

	sut := usecases.NewDefaultGetUserConfig(mockRepository)

	expectedResult := domain.Config{
		LastOpenedProjectId: "test-project-id",
		EnvironmentPaths: []domain.EnvironmentPath{
			{
				Id:   "test-env-path-id",
				Path: "test/path",
			},
		},
	}
	mockRepository.On("GetOrCreate").Return(&expectedResult, nil)

	config, err := sut.Execute()

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, *config)
}
