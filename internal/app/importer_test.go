package app_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	"gomander/internal/testutils/mocks"
)

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

func TestApp_GetProjectToImport(t *testing.T) {
	t.Run("Should return project import", func(t *testing.T) {
		// Arrange
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade: mockRuntimeFacade,
			FsFacade:      mockFsFacade,
		})

		basicProjectJson := app.ProjectExportJSONv1{
			Version:       1,
			Name:          "Name",
			Commands:      make([]app.CommandJSONv1, 0),
			CommandGroups: make([]app.CommandGroupJSONv1, 0),
		}

		basicProjectJsonBytes, err := json.Marshal(basicProjectJson)
		assert.NoError(t, err)

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("/path/to/project.json", nil)
		mockFsFacade.On("ReadFile", "/path/to/project.json").Return(basicProjectJsonBytes, nil)

		// Act
		toImport, err := a.GetProjectToImport()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &basicProjectJson, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return error if there is a problem opening the file dialog", func(t *testing.T) {
		// Arrange
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade: mockRuntimeFacade,
			FsFacade:      mockFsFacade,
		})

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("", assert.AnError)

		// Act
		toImport, err := a.GetProjectToImport()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return nil if the user cancels the file dialog", func(t *testing.T) {
		// Arrange
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade: mockRuntimeFacade,
			FsFacade:      mockFsFacade,
		})

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("", nil)

		// Act
		toImport, err := a.GetProjectToImport()

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return error if there is a problem reading the file", func(t *testing.T) {
		// Arrange
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade: mockRuntimeFacade,
			FsFacade:      mockFsFacade,
		})

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("/path/to/project.json", nil)
		mockFsFacade.On("ReadFile", "/path/to/project.json").Return([]byte{}, assert.AnError)

		// Act
		toImport, err := a.GetProjectToImport()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
}
