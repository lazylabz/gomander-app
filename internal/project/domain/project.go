package domain

type Project struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	WorkingDirectory string   `json:"workingDirectory"`
	FailurePatterns  []string `json:"failurePatterns,omitempty"`
}
