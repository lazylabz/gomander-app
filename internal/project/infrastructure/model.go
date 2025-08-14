package infrastructure

type ProjectModel struct {
	Id               string `gorm:"primaryKey;column:id"`
	Name             string `gorm:"column:name"`
	WorkingDirectory string `gorm:"column:working_directory"`
}

func (ProjectModel) TableName() string {
	return "project"
}
