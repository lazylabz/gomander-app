package fs_test

import (
	"errors"
	stdruntime "runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/dialog"
	dialogtest "gomander/internal/dialog/test"
	"gomander/internal/facade/test"
	"gomander/internal/uihelpers/fs"
)

func TestUIFsHelper_NewUIFsHelper(t *testing.T) {
	// Arrange
	mockDialogs := new(dialogtest.MockDialogs)
	mockRuntime := new(test.MockRuntimeFacade)

	// Act
	helper := fs.NewUIFsHelper(mockDialogs, mockRuntime)

	// Assert
	assert.NotNil(t, helper)
}

func TestUIFsHelper_AskForDirPath(t *testing.T) {
	t.Run("Should return directory path when successful", func(t *testing.T) {
		// Arrange
		mockDialogs := new(dialogtest.MockDialogs)
		mockRuntime := new(test.MockRuntimeFacade)
		helper := fs.NewUIFsHelper(mockDialogs, mockRuntime)

		expectedPath := "/some/directory/path"
		mockDialogs.On("AskForDirectory", dialog.PickDirectoryRequest{}).Return(expectedPath, nil)

		// Act
		path, err := helper.AskForDirPath()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, path)
		mockDialogs.AssertExpectations(t)
	})

	t.Run("Should return error when asking for a directory fails", func(t *testing.T) {
		// Arrange
		mockDialogs := new(dialogtest.MockDialogs)
		mockRuntime := new(test.MockRuntimeFacade)
		helper := fs.NewUIFsHelper(mockDialogs, mockRuntime)

		expectedError := errors.New("dialog error")
		mockDialogs.On("AskForDirectory", dialog.PickDirectoryRequest{}).Return("", expectedError)

		// Act
		path, err := helper.AskForDirPath()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, path)
		mockDialogs.AssertExpectations(t)
	})
}

func TestUIFsHelper_OpenFileFolder(t *testing.T) {
	// Arrange
	mockDialogs := new(dialogtest.MockDialogs)
	mockRuntime := new(test.MockRuntimeFacade)
	helper := fs.NewUIFsHelper(mockDialogs, mockRuntime)

	filePath := "/some/directory/file.txt"
	expectedFolderPath := "/some/directory"

	if stdruntime.GOOS == "windows" {
		filePath = "C:\\some\\directory\\file.txt"
		expectedFolderPath = "C:\\some\\directory"
	}

	mockRuntime.On("OpenFolderInFileManager", expectedFolderPath).Return(nil)

	// Act
	err := helper.OpenFileFolder(filePath)
	assert.NoError(t, err)

	// Assert
	mockRuntime.AssertExpectations(t)
}
