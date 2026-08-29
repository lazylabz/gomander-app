//go:build darwin

package infrastructure

import (
	"fmt"
	"runtime"
)

func binaryName() string {
	return fmt.Sprintf("gomander-darwin-%s.dmg", runtime.GOARCH)
}

func pathToOpen(binaryPath string) string {
	return binaryPath
}
