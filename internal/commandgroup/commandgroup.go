package commandgroup

type CommandGroup struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	CommandIds []string `json:"commands"`
}
