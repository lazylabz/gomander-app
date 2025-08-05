package project

type Project struct {
	Id                   string             `json:"id"`
	Name                 string             `json:"name"`
	BaseWorkingDirectory string             `json:"baseWorkingDirectory"`
	Commands             map[string]Command `json:"commands"`
	CommandGroups        []CommandGroup     `json:"commandGroups"`
}

type ExportableProject struct {
	Id            string             `json:"id"`
	Name          string             `json:"name"`
	Commands      map[string]Command `json:"commands"`
	CommandGroups []CommandGroup     `json:"commandGroups"`
}

func (p *Project) ToExportable() *ExportableProject {
	return &ExportableProject{
		Id:            p.Id,
		Name:          p.Name,
		Commands:      p.Commands,
		CommandGroups: p.CommandGroups,
	}
}

func (e *ExportableProject) ToProject(baseWorkingDirectory string) *Project {
	return &Project{
		Id:                   e.Id,
		Name:                 e.Name,
		BaseWorkingDirectory: baseWorkingDirectory,
		Commands:             e.Commands,
		CommandGroups:        e.CommandGroups,
	}
}
