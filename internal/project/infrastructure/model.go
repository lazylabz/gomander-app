package infrastructure

import "github.com/lib/pq"

type ProjectModel struct {
	Id               string         `gorm:"primaryKey;column:id"`
	Name             string         `gorm:"column:name"`
	WorkingDirectory string         `gorm:"column:working_directory"`
	FailurePatterns  pq.StringArray `gorm:"type:text[]"`
}

func (ProjectModel) TableName() string {
	return "project"
}
