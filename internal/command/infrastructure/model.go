package commandinfrastructure

type CommandModel struct {
	Id               string `gorm:"primaryKey;column:id"`
	ProjectId        string `gorm:"column:project_id"`
	Name             string `gorm:"column:name"`
	Command          string `gorm:"column:command"`
	WorkingDirectory string `gorm:"column:working_directory"`
	Position         int    `gorm:"column:position"`
}

func (CommandModel) TableName() string {
	return "command"
}
