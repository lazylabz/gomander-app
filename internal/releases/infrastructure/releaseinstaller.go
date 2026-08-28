package infrastructure

import (
	"gomander/internal/facade"
)

// OSReleaseInstaller installs a downloaded Release by handing its binary to the
// operating system, which opens it with whatever it uses for that kind of file.
type OSReleaseInstaller struct {
	osFacade   facade.OSFacade
	openFacade facade.OpenFacade
}

func NewOSReleaseInstaller(osFacade facade.OSFacade, openFacade facade.OpenFacade) *OSReleaseInstaller {
	return &OSReleaseInstaller{
		osFacade:   osFacade,
		openFacade: openFacade,
	}
}

func (i *OSReleaseInstaller) Install(binaryPath string) error {
	if _, err := i.osFacade.Stat(binaryPath); err != nil {
		return err
	}

	return i.openFacade.Run(pathToOpen(binaryPath))
}
