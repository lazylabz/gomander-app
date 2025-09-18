package usecases_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/command/application/usecases"
	"gomander/internal/runner/test"
)

func TestDefaultGetRunningCommandIds_Execute(t *testing.T) {
	t.Run("Should return empty list when there are no running commands", func(t *testing.T) {
		// Arrange
		mockRunner := new(test.MockRunner)
		sut := usecases.NewGetRunningCommandIds(mockRunner)

		mockRunner.On("GetRunningCommandIds").Return([]string{})

		// Act
		result := sut.Execute()

		// Assert
		assert.Empty(t, result)
		mockRunner.AssertExpectations(t)
	})

	t.Run("Should return list of running command ids", func(t *testing.T) {
		// Arrange
		mockRunner := new(test.MockRunner)
		sut := usecases.NewGetRunningCommandIds(mockRunner)

		expectedIds := []string{"cmd-1", "cmd-2", "cmd-3"}
		mockRunner.On("GetRunningCommandIds").Return(expectedIds)

		// Act
		result := sut.Execute()

		// Assert
		assert.Equal(t, expectedIds, result)
		assert.Len(t, result, len(expectedIds))
		mockRunner.AssertExpectations(t)
	})
}
