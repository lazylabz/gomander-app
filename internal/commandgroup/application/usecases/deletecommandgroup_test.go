package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/commandgroup/application/usecases"
	"gomander/internal/commandgroup/domain/test"
	"gomander/internal/event"
	test2 "gomander/internal/event/test"
)

func TestDefaultDeleteCommandGroup_Execute(t *testing.T) {
	t.Run("Should delete a command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(test.MockCommandGroupRepository)
		mockEventEventEmitter := new(test2.MockEventEmitter)

		projectId := "project1"

		sut := usecases.NewDeleteCommandGroup(mockCommandGroupRepository, mockEventEventEmitter)

		paramCommandGroup := test.NewCommandGroupBuilder().
			WithProjectId(projectId).
			Build()

		mockCommandGroupRepository.On("Delete", paramCommandGroup.Id).Return(nil)
		mockEventEventEmitter.On("EmitEvent", event.CommandGroupDeleted, paramCommandGroup.Id).Return()

		// Act
		err := sut.Execute(paramCommandGroup.Id)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockEventEventEmitter,
		)
	})

	t.Run("Should return an error if failing to delete the command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(test.MockCommandGroupRepository)
		mockEventEventEmitter := new(test2.MockEventEmitter)

		projectId := "project1"

		sut := usecases.NewDeleteCommandGroup(mockCommandGroupRepository, mockEventEventEmitter)

		paramCommandGroup := test.NewCommandGroupBuilder().
			WithProjectId(projectId).
			Build()

		mockCommandGroupRepository.On("Delete", paramCommandGroup.Id).Return(errors.New("failed to delete command group"))

		// Act
		err := sut.Execute(paramCommandGroup.Id)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockEventEventEmitter,
		)
	})
}
