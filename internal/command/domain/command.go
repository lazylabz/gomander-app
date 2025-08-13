package domain

type Command struct {
	Id               string `json:"id"`
	ProjectId        string `json:"projectId"`
	Name             string `json:"name"`
	Command          string `json:"command"`
	WorkingDirectory string `json:"workingDirectory"`
	Position         int    `json:"position"`
}

func (c *Command) Equals(other *Command) bool {
	if c == nil || other == nil {
		return false
	}
	if c.Id != other.Id || c.ProjectId != other.ProjectId || c.Name != other.Name || c.Command != other.Command || c.WorkingDirectory != other.WorkingDirectory || c.Position != other.Position {
		return false
	}
	return true
}
