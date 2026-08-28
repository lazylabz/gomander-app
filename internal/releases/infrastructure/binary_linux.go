//go:build linux

package infrastructure

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func binaryName() string {
	return fmt.Sprintf("gomander-linux-%s", runtime.GOARCH)
}

// The Linux release is the executable itself rather than an installer, so the
// user is shown the folder it landed in instead of having it launched for them.
func pathToOpen(binaryPath string) string {
	return filepath.Dir(binaryPath)
}
