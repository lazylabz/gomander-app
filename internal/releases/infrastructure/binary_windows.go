//go:build windows

package infrastructure

import (
	"fmt"
	"runtime"
)

func binaryName() string {
	return fmt.Sprintf("gomander-windows-%s-installer.exe", runtime.GOARCH)
}

func pathToOpen(binaryPath string) string {
	return binaryPath
}
