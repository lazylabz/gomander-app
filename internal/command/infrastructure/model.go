package infrastructure

type CommandModel struct {
	Id               string `gorm:"primaryKey;column:id"`
	ProjectId        string `gorm:"column:project_id"`
	Name             string `gorm:"column:name"`
	Command          string `gorm:"column:command"`
	WorkingDirectory string `gorm:"column:working_directory"`
	Position         int    `gorm:"column:position"`
	Link             string `gorm:"column:link"`
	ErrorPatterns    string `gorm:"column:error_patterns"`
}

func (CommandModel) TableName() string {
	return "command"
}
