//go:build windows

package releases

import (
	"fmt"
	"runtime"
)

func (rh *ReleaseHelper) runBinary(binaryPath string) error {
	return rh.openFacade.Run(binaryPath)
}

func (rh *ReleaseHelper) getBinaryName(_ string) string {
	return fmt.Sprintf("gomander-windows-%s-installer.exe", runtime.GOARCH)
}
