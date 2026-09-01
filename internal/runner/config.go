package runner

// HostEnvironment identifies an operating-system environment whose process
// behavior has been explicitly validated by Gomander.
type HostEnvironment string

const (
	HostEnvironmentWindows10 HostEnvironment = "windows10"
)

// Config controls compatibility-sensitive runner behavior.
type Config struct {
	ConPTYEnvironments []HostEnvironment
}

func (c Config) enablesConPTY(environment HostEnvironment) bool {
	for _, allowed := range c.ConPTYEnvironments {
		if allowed == environment {
			return true
		}
	}
	return false
}
