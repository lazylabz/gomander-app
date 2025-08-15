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
