package usecases_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/config/application/usecases"
	"gomander/internal/config/domain"
)

func TestDefaultGetUserConfig_Execute(t *testing.T) {
	mockRepository := new(MockUserConfigRepository)

	sut := usecases.NewGetUserConfig(mockRepository)

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
