package os_internal

import (
	"runtime"
)

type UIOsHelper struct {
}

func NewUIOsHelper() *UIOsHelper {
	return &UIOsHelper{}
}

func (uh *UIOsHelper) GetOs() string {
	return runtime.GOOS
}
