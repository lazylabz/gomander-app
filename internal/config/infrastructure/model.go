package infrastructure

type ConfigModel struct {
	Id                  int    `gorm:"primaryKey;column:id"`
	LastOpenedProjectId string `gorm:"column:last_opened_project_id"`
	LogLineLimit        int    `gorm:"column:log_line_limit"`
}

func (ConfigModel) TableName() string {
	return "user_config"
}

type EnvironmentPathModel struct {
	Id   string `gorm:"primaryKey;column:id"`
	Path string `gorm:"column:path"`
}

func (EnvironmentPathModel) TableName() string {
	return "user_config_environment_paths"
}
