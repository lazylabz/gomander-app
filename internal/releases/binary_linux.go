//go:build linux

package releases

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func (rh *ReleaseHelper) runBinary(binaryPath string) error {
	return rh.openFacade.Run(filepath.Dir(binaryPath))
}

func (rh *ReleaseHelper) getBinaryName(_ string) string {
	return fmt.Sprintf("gomander-linux-%s", runtime.GOARCH)
}
