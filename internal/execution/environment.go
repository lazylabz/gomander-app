// Package execution holds the environment a Command runs in, resolved before
// the Runner is involved, so the Runner is told where to run and with which
// paths instead of working it out from what it would have to look up.
package execution

type Environment struct {
	Paths                []string
	BaseWorkingDirectory string
}
