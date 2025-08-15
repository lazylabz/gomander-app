package domain

import (
	"gomander/internal/command/domain"
)

type CommandGroup struct {
	Id        string           `json:"id"`
	ProjectId string           `json:"projectId"`
	Name      string           `json:"name"`
	Commands  []domain.Command `json:"commands"`
	Position  int              `json:"position"`
}
