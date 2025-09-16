package fs_test

import (
	"context"
	"errors"
	stdruntime "runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gomander/internal/facade/test"
	"gomander/internal/uihelpers/fs"
)

func TestUIFsHelper_NewUIFsHelper(t *testing.T) {
	// Arrange
	mockRuntime := new(test.MockRuntimeFacade)

	// Act
	helper := fs.NewUIFsHelper(mockRuntime)

	// Assert
	assert.NotNil(t, helper)
}

func TestUIFsHelper_SetContext(t *testing.T) {
	// Arrange
	mockRuntime := new(test.MockRuntimeFacade)
	helper := fs.NewUIFsHelper(mockRuntime)
	ctx := context.Background()

	// Act
	helper.SetContext(ctx)

	// No Assert needed as this is a simple setter method
	// The method doesn't return anything and just sets an internal field
}

func TestUIFsHelper_AskForDirPath(t *testing.T) {
	t.Run("Should return directory path when successful", func(t *testing.T) {
		// Arrange
		mockRuntime := new(test.MockRuntimeFacade)
		helper := fs.NewUIFsHelper(mockRuntime)
		ctx := context.Background()
		helper.SetContext(ctx)

		expectedPath := "/some/directory/path"
		mockRuntime.On("OpenDirectoryDialog", ctx, runtime.OpenDialogOptions{}).Return(expectedPath, nil)

		// Act
		path, err := helper.AskForDirPath()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, path)
		mockRuntime.AssertExpectations(t)
	})

	t.Run("Should return error when OpenDirectoryDialog fails", func(t *testing.T) {
		// Arrange
		mockRuntime := new(test.MockRuntimeFacade)
		helper := fs.NewUIFsHelper(mockRuntime)
		ctx := context.Background()
		helper.SetContext(ctx)

		expectedError := errors.New("dialog error")
		mockRuntime.On("OpenDirectoryDialog", ctx, runtime.OpenDialogOptions{}).Return("", expectedError)

		// Act
		path, err := helper.AskForDirPath()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, path)
		mockRuntime.AssertExpectations(t)
	})
}

func TestUIFsHelper_OpenDirectoryDialog(t *testing.T) {
	// Arrange
	mockRuntime := new(test.MockRuntimeFacade)
	helper := fs.NewUIFsHelper(mockRuntime)
	ctx := context.Background()
	helper.SetContext(ctx)

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
