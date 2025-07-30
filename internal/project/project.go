package project

type Project struct {
	Id            string             `json:"id"`
	Name          string             `json:"name"`
	Commands      map[string]Command `json:"commands"`
	CommandGroups []CommandGroup     `json:"command_groups"`
}
