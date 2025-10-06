package usecases_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/facade/test"
	"gomander/internal/project/application/usecases"
	projectdomain "gomander/internal/project/domain"
)

var testCases = []struct {
	fileType     usecases.FileType
	dialogPath   string
	fileData     interface{}
	expectedData *projectdomain.ProjectExportJSONv1
}{
	{
		fileType:   usecases.FileTypeGomander,
		dialogPath: "/path/to/gomander_project.json",
		fileData: projectdomain.ProjectExportJSONv1{
			Version:       1,
			Name:          "Gomander Project",
			Commands:      make([]projectdomain.CommandJSONv1, 0),
			CommandGroups: make([]projectdomain.CommandGroupJSONv1, 0),
		},
		expectedData: &projectdomain.ProjectExportJSONv1{
			Version:       1,
			Name:          "Gomander Project",
			Commands:      make([]projectdomain.CommandJSONv1, 0),
			CommandGroups: make([]projectdomain.CommandGroupJSONv1, 0),
		},
	},
	{
		fileType:   usecases.FileTypePackageJSON,
		dialogPath: "/path/to/package.json",
		fileData: map[string]interface{}{
			"name": "My NPM Project",
			"scripts": map[string]interface{}{
				"start": "node index.js",
			},
		},
		expectedData: &projectdomain.ProjectExportJSONv1{
			Version:          1,
			Name:             "My NPM Project",
			WorkingDirectory: "/path/to",
			Commands: []projectdomain.CommandJSONv1{
				{
					Id:               "cmd-1",
					Name:             "start",
					Command:          "node index.js",
					WorkingDirectory: "",
				},
			},
		},
	},
}

func TestDefaultGetProjectToImport_Execute(t *testing.T) {
	for _, testCase := range testCases {
		t.Run("Should return project import for "+string(testCase.fileType), func(t *testing.T) {
			// Arrange

			mockRuntimeFacade := new(test.MockRuntimeFacade)
			mockFsFacade := new(test.MockFsFacade)

			sut := usecases.NewGetProjectToImport(context.Background(), mockRuntimeFacade, mockFsFacade)

			dataBytes, err := json.Marshal(testCase.fileData)
			assert.NoError(t, err)

			mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return(testCase.dialogPath, nil)
			mockFsFacade.On("ReadFile", testCase.dialogPath).Return(dataBytes, nil)

			// Act
			toImport, err := sut.Execute(testCase.fileType)

			// Assert
			assert.NoError(t, err)

			// Assert keys individually to avoid checking random ids
			assert.Equal(t, testCase.expectedData.Name, toImport.Name)
			assert.Equal(t, testCase.expectedData.WorkingDirectory, toImport.WorkingDirectory)
			assert.Equal(t, len(testCase.expectedData.Commands), len(toImport.Commands))
			if len(testCase.expectedData.Commands) > 0 {
				assert.NotNil(t, toImport.Commands[0].Id)
				assert.Equal(t, testCase.expectedData.Commands[0].Name, toImport.Commands[0].Name)
				assert.Equal(t, testCase.expectedData.Commands[0].Command, toImport.Commands[0].Command)
				assert.Equal(t, testCase.expectedData.Commands[0].WorkingDirectory, toImport.Commands[0].WorkingDirectory)
			}
			assert.Equal(t, len(testCase.expectedData.CommandGroups), len(toImport.CommandGroups))

			mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
		})
	}

	t.Run("Should return error if there is a problem opening the file dialog", func(t *testing.T) {
		// Arrange

		mockRuntimeFacade := new(test.MockRuntimeFacade)
		mockFsFacade := new(test.MockFsFacade)

		sut := usecases.NewGetProjectToImport(context.Background(), mockRuntimeFacade, mockFsFacade)

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("", assert.AnError)

		// Act
		toImport, err := sut.Execute(usecases.FileTypeGomander)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return nil if the user cancels the file dialog", func(t *testing.T) {
		// Arrange

		mockRuntimeFacade := new(test.MockRuntimeFacade)
		mockFsFacade := new(test.MockFsFacade)

		sut := usecases.NewGetProjectToImport(context.Background(), mockRuntimeFacade, mockFsFacade)

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("", nil)

		// Act
		toImport, err := sut.Execute(usecases.FileTypeGomander)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return error if there is a problem reading the file", func(t *testing.T) {
		// Arrange

		mockRuntimeFacade := new(test.MockRuntimeFacade)
		mockFsFacade := new(test.MockFsFacade)

		sut := usecases.NewGetProjectToImport(context.Background(), mockRuntimeFacade, mockFsFacade)

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("/path/to/project.json", nil)
		mockFsFacade.On("ReadFile", "/path/to/project.json").Return([]byte{}, assert.AnError)

		// Act
		toImport, err := sut.Execute(usecases.FileTypeGomander)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
}
