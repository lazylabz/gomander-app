package fs_test

import (
	"context"
	"errors"
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

func TestUIFsHelper_GetDirPath(t *testing.T) {
	t.Run("Should return directory path when successful", func(t *testing.T) {
		// Arrange
		mockRuntime := new(test.MockRuntimeFacade)
		helper := fs.NewUIFsHelper(mockRuntime)
		ctx := context.Background()
		helper.SetContext(ctx)

		expectedPath := "/some/directory/path"
		mockRuntime.On("OpenDirectoryDialog", ctx, runtime.OpenDialogOptions{}).Return(expectedPath, nil)

		// Act
		path, err := helper.GetDirPath()

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
		path, err := helper.GetDirPath()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, path)
		mockRuntime.AssertExpectations(t)
	})
}
