package usecases

import (
	"gomander/internal/releases"
)

type InstallReleaseAndQuit struct {
	releaseInstaller releases.ReleaseInstaller
	appControl       releases.AppControl
}

func NewInstallReleaseAndQuit(releaseInstaller releases.ReleaseInstaller, appControl releases.AppControl) *InstallReleaseAndQuit {
	return &InstallReleaseAndQuit{
		releaseInstaller: releaseInstaller,
		appControl:       appControl,
	}
}

// Execute quits only once the operating system has taken the binary: a failed
// install leaves the user with the app they already had.
func (uc *InstallReleaseAndQuit) Execute(binaryPath string) error {
	if err := uc.releaseInstaller.Install(binaryPath); err != nil {
		return err
	}

	uc.appControl.CloseApp()

	return nil
}
