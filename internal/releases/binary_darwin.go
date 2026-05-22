//go:build darwin

package releases

import (
	"fmt"
	"runtime"
)

func (rh *ReleaseHelper) runBinary(binaryPath string) error {
	return rh.openFacade.Run(binaryPath)
}

func (rh *ReleaseHelper) getBinaryName(_ string) string {
	return fmt.Sprintf("gomander-darwin-%s.dmg", runtime.GOARCH)
}
