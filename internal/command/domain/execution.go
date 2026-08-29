package domain

import (
	"path/filepath"
	"strings"
)

// MatchesErrorPattern answers whether a line of this Command's output marks it
// as failing. An Error Pattern is a plain substring, not a regex, and the first
// one the line contains settles the answer.
func (c Command) MatchesErrorPattern(line string) bool {
	for _, errorPattern := range c.ErrorPatterns {
		if strings.Contains(line, errorPattern) {
			return true
		}
	}

	return false
}

// ResolveWorkingDirectory answers the directory this Command runs in, given the
// base working directory of the Execution Environment: no directory of its own
// means the base, an absolute one wins over the base, and a relative one hangs
// off it.
func (c Command) ResolveWorkingDirectory(baseWorkingDirectory string) string {
	if c.WorkingDirectory == "" {
		return baseWorkingDirectory
	}

	if filepath.IsAbs(c.WorkingDirectory) {
		return c.WorkingDirectory
	}

	return filepath.Join(baseWorkingDirectory, c.WorkingDirectory)
}
