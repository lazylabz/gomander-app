package project

type Project struct {
	Id                   string             `json:"id"`
	Name                 string             `json:"name"`
	BaseWorkingDirectory string             `json:"baseWorkingDirectory"`
	Commands             map[string]Command `json:"commands"`
	OrderedCommandIds    []string           `json:"orderedCommandIds"`
	CommandGroups        []CommandGroup     `json:"commandGroups"`
}
