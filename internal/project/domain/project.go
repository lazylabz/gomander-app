package domain

type Project struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	WorkingDirectory string `json:"workingDirectory"`
}

func (p *Project) Equals(other *Project) bool {
	if p == nil || other == nil {
		return p == other
	}
	return p.Id == other.Id &&
		p.Name == other.Name &&
		p.WorkingDirectory == other.WorkingDirectory
}
