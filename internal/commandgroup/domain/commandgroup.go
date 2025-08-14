package domain

import (
	"gomander/internal/command/domain"
)

type CommandGroup struct {
	Id        string
	ProjectId string
	Name      string
	Commands  []domain.Command
	Position  int
}

func (cg *CommandGroup) Equals(other *CommandGroup) bool {
	if cg == nil || other == nil {
		return false
	}
	if cg.Id != other.Id || cg.ProjectId != other.ProjectId || cg.Name != other.Name || cg.Position != other.Position {
		return false
	}
	if len(cg.Commands) != len(other.Commands) {
		return false
	}
	for i, cmd := range cg.Commands {
		if !cmd.Equals(&other.Commands[i]) {
			return false
		}
	}
	return true
}
