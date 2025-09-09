package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/domain/test"
	"gomander/internal/commandgroup/application/usecases"
	"gomander/internal/commandgroup/domain"
	test2 "gomander/internal/commandgroup/domain/test"
	configdomain "gomander/internal/config/domain"
	test3 "gomander/internal/config/domain/test"
)

func TestDefaultCreateCommandGroup_Execute(t *testing.T) {
	t.Run("Should create a command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)

		projectId := "project1"

		sut := usecases.NewCreateCommandGroup(mockUserConfigRepository, mockCommandGroupRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		otherCommandGroup := test2.NewCommandGroupBuilder().
			WithProjectId(projectId).
			Build()

		commandGroupBuilder := test2.NewCommandGroupBuilder().
			WithProjectId(projectId).
			WithPosition(0).
			WithCommands(
				test.NewCommandBuilder().
					WithProjectId(projectId).
					Build(),
			)

		paramCommandGroup := commandGroupBuilder.Build()
		expectedCommandGroupCall := commandGroupBuilder.WithPosition(1).Build()

		mockCommandGroupRepository.On("GetAll", projectId).Return([]domain.CommandGroup{otherCommandGroup}, nil)
		mockCommandGroupRepository.On("Create", &expectedCommandGroupCall).Return(nil)

		// Act
		err := sut.Execute(&paramCommandGroup)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if failing to retrieve user config", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)

		sut := usecases.NewCreateCommandGroup(mockUserConfigRepository, mockCommandGroupRepository)

		expectedError := errors.New("failed to get user config")
		mockUserConfigRepository.On("GetOrCreate").Return(nil, expectedError)

		paramCommandGroup := test2.NewCommandGroupBuilder().
			WithProjectId("project1").
			Build()

		// Act
		err := sut.Execute(&paramCommandGroup)

		// Assert
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if failing to retrieve existing command groups", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)

		projectId := "project1"

		sut := usecases.NewCreateCommandGroup(mockUserConfigRepository, mockCommandGroupRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		paramCommandGroup := test2.NewCommandGroupBuilder().
			WithProjectId(projectId).
			Build()

		mockCommandGroupRepository.On("GetAll", projectId).Return(make([]domain.CommandGroup, 0), errors.New("failed to get command groups"))

		// Act
		err := sut.Execute(&paramCommandGroup)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if failing to save the command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)

		projectId := "project1"

		sut := usecases.NewCreateCommandGroup(mockUserConfigRepository, mockCommandGroupRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		paramCommandGroup := test2.NewCommandGroupBuilder().
			WithProjectId(projectId).
			Build()

		mockCommandGroupRepository.On("GetAll", projectId).Return(make([]domain.CommandGroup, 0), nil)
		mockCommandGroupRepository.On("Create", &paramCommandGroup).Return(errors.New("failed to create command group"))

		// Act
		err := sut.Execute(&paramCommandGroup)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})
}
