package infrastructure_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	facadetest "gomander/internal/facade/test"
	"gomander/internal/releases/infrastructure"
)

func TestOSReleaseInstaller_Install(t *testing.T) {
	binaryPath := "/some/path/to/binary"

	t.Run("Should hand the binary to the operating system", func(t *testing.T) {
		// Arrange
		mockOSFacade := new(facadetest.MockOSFacade)
		mockOpenFacade := new(facadetest.MockOpenFacade)

		mockOSFacade.On("Stat", binaryPath).Return(nil, nil)
		mockOpenFacade.On("Run", expectedOpenArg(binaryPath)).Return(nil)

		installer := infrastructure.NewOSReleaseInstaller(mockOSFacade, mockOpenFacade)

		// Act
		err := installer.Install(binaryPath)

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockOSFacade, mockOpenFacade)
	})

	t.Run("Should return an error when the binary does not exist", func(t *testing.T) {
		// Arrange
		mockOSFacade := new(facadetest.MockOSFacade)
		mockOpenFacade := new(facadetest.MockOpenFacade)

		mockOSFacade.On("Stat", binaryPath).Return(nil, os.ErrNotExist)

		installer := infrastructure.NewOSReleaseInstaller(mockOSFacade, mockOpenFacade)

		// Act
		err := installer.Install(binaryPath)

		// Assert
		assert.ErrorIs(t, err, os.ErrNotExist)
		mockOpenFacade.AssertNotCalled(t, "Run", mock.Anything)
		mock.AssertExpectationsForObjects(t, mockOSFacade)
	})

	t.Run("Should return the error when the operating system refuses the binary", func(t *testing.T) {
		// Arrange
		mockOSFacade := new(facadetest.MockOSFacade)
		mockOpenFacade := new(facadetest.MockOpenFacade)

		mockOSFacade.On("Stat", binaryPath).Return(nil, nil)
		mockOpenFacade.On("Run", expectedOpenArg(binaryPath)).Return(assert.AnError)

		installer := infrastructure.NewOSReleaseInstaller(mockOSFacade, mockOpenFacade)

		// Act
		err := installer.Install(binaryPath)

		// Assert
		assert.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, mockOSFacade, mockOpenFacade)
	})
}

// expectedOpenArg is the Linux exception spelled out: there the release is the
// executable itself, so the user is shown the folder rather than having it run.
func expectedOpenArg(binaryPath string) string {
	if runtime.GOOS == "linux" {
		return filepath.Dir(binaryPath)
	}

	return binaryPath
}
