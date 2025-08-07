package app

type EditProjectDTO struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	BaseWorkingDirectory string `json:"baseWorkingDirectory"`
}
