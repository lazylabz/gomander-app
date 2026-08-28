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

// CommandGroupWithCommandIds is a Command Group that names the Commands it
// holds instead of carrying copies of them, in the order it holds them. It
// exists so callers can move off the embedded form one at a time; once none is
// left it becomes the only shape of a Command Group.
type CommandGroupWithCommandIds struct {
	Id         string
	ProjectId  string
	Name       string
	CommandIds []string
	Position   int
}
