package infrastructure

type CommandGroupModel struct {
	Id        string `gorm:"primaryKey;column:id"`
	ProjectId string `gorm:"column:project_id"`
	Name      string `gorm:"column:name"`
	Position  int    `gorm:"column:position"`
}

func (CommandGroupModel) TableName() string {
	return "command_group"
}

type CommandToCommandGroupModel struct {
	CommandGroupId string `gorm:"primaryKey;column:command_group_id"`
	CommandId      string `gorm:"primaryKey;column:command_id"`
	Position       int    `gorm:"column:position"`
}

func (CommandToCommandGroupModel) TableName() string {
	return "command_group_command"
}
