package infrastructure

type ConfigModel struct {
	Id                  int     `gorm:"primaryKey;column:id"`
	LastOpenedProjectId string  `gorm:"column:last_opened_project_id"`
	Locale              string  `gorm:"column:locale"`
	LogsFontSize        int     `gorm:"column:logs_font_size;default:15"`
	ZoomLevel           float64 `gorm:"column:zoom_level;default:1"`
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
