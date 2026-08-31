// Package releases is how Gomander updates itself: it knows which Release it is
// running, asks the release feed which one is published, and downloads and
// installs a newer one. Reaching the network, writing a binary to disk and
// quitting the app all happen behind the ports declared here, so the flow is
// driven like every other backend operation and none of it needs a desktop
// toolkit to be exercised.
package releases

// CurrentRelease is the Release this binary is.
const CurrentRelease = "v1.9.0"

// ReleaseFeed is where the app learns which Releases have been published. It
// answers with the latest one's version, or with an empty string when the feed
// lists none - a feed it cannot make sense of is an error, not an empty one.
type ReleaseFeed interface {
	GetLatestRelease() (version string, err error)
}

// ReleaseDownloader writes a published Release's binary where the operating
// system can reach it, and answers with the path it wrote.
type ReleaseDownloader interface {
	Download(version string) (binaryPath string, err error)
}

// ReleaseInstaller hands a downloaded Release's binary to the operating system,
// which is what installing means on each platform.
type ReleaseInstaller interface {
	Install(binaryPath string) error
}

// AppControl is the app's own lifecycle, as far as the update flow needs it:
// installing a downloaded Release means quitting so the new binary can take
// over.
type AppControl interface {
	CloseApp()
}
