package domain

import "strings"

type Command struct {
	Id               string `json:"id"`
	ProjectId        string `json:"projectId"`
	Name             string `json:"name"`
	Command          string `json:"command"`
	WorkingDirectory string `json:"workingDirectory"`
	Position         int    `json:"position"`
	Link             string `json:"link"`
	ErrorPatterns    string `json:"errorPatterns"`
}

// GetErrorPatterns splits the ErrorPatterns string (delimited by newlines) into a slice of patterns.
func (c *Command) GetErrorPatterns() []string {
	var patterns []string

	for _, pattern := range strings.Split(c.ErrorPatterns, "\n") {
		trimmed := strings.TrimSpace(pattern)
		if trimmed != "" {
			patterns = append(patterns, trimmed)
		}
	}

	return patterns
}
