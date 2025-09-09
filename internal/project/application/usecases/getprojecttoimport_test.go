package usecases_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/project/application/usecases"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/testutils/mocks"
)

func TestDefaultGetProjectToImport_Execute(t *testing.T) {
	t.Run("Should return project import", func(t *testing.T) {
		// Arrange

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)

		sut := usecases.NewGetProjectToImport(context.Background(), mockRuntimeFacade, mockFsFacade)

		basicProjectJson := projectdomain.ProjectExportJSONv1{
			Version:       1,
			Name:          "Name",
			Commands:      make([]projectdomain.CommandJSONv1, 0),
			CommandGroups: make([]projectdomain.CommandGroupJSONv1, 0),
		}

		basicProjectJsonBytes, err := json.Marshal(basicProjectJson)
		assert.NoError(t, err)

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("/path/to/project.json", nil)
		mockFsFacade.On("ReadFile", "/path/to/project.json").Return(basicProjectJsonBytes, nil)

		// Act
		toImport, err := sut.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &basicProjectJson, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return error if there is a problem opening the file dialog", func(t *testing.T) {
		// Arrange

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)

		sut := usecases.NewGetProjectToImport(context.Background(), mockRuntimeFacade, mockFsFacade)

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("", assert.AnError)

		// Act
		toImport, err := sut.Execute()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return nil if the user cancels the file dialog", func(t *testing.T) {
		// Arrange

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)

		sut := usecases.NewGetProjectToImport(context.Background(), mockRuntimeFacade, mockFsFacade)

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("", nil)

		// Act
		toImport, err := sut.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return error if there is a problem reading the file", func(t *testing.T) {
		// Arrange

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)

		sut := usecases.NewGetProjectToImport(context.Background(), mockRuntimeFacade, mockFsFacade)

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("/path/to/project.json", nil)
		mockFsFacade.On("ReadFile", "/path/to/project.json").Return([]byte{}, assert.AnError)

		// Act
		toImport, err := sut.Execute()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
}
