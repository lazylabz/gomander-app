package runner

import (
	"io"
	"os"
	"strings"
)

type commandProcess interface {
	Start() ([]io.ReadCloser, error)
	Wait() error
	PID() int
}

func commandEnvironment(environmentPaths []string) []string {
	environment := append(os.Environ(), "FORCE_COLOR=1", "TERM=xterm-256color")
	if len(environmentPaths) == 0 {
		return environment
	}

	newPath := strings.Join(environmentPaths, string(os.PathListSeparator)) +
		string(os.PathListSeparator) + os.Getenv("PATH")

	for i, variable := range environment {
		if strings.HasPrefix(strings.ToUpper(variable), "PATH=") {
			environment[i] = "PATH=" + newPath
			return environment
		}
	}

	return append(environment, "PATH="+newPath)
}
