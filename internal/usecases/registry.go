package usecases

import (
	commandusecases "gomander/internal/command/application/usecases"
	commandgroupusecases "gomander/internal/commandgroup/application/usecases"
	configusecases "gomander/internal/config/application/usecases"
	localizationusecases "gomander/internal/localization/application/usecases"
	projectusecases "gomander/internal/project/application/usecases"
)

type Registry struct {
	// Configuration
	GetUserConfig  *configusecases.GetUserConfig
	SaveUserConfig *configusecases.SaveUserConfig
	// Localization
	GetTranslation        *localizationusecases.GetTranslation
	GetSupportedLanguages *localizationusecases.GetSupportedLanguages
	// Projects
	GetCurrentProject    *projectusecases.GetCurrentProject
	GetAvailableProjects *projectusecases.GetAvailableProjects
	OpenProject          *projectusecases.OpenProject
	CreateProject        *projectusecases.CreateProject
	EditProject          *projectusecases.EditProject
	CloseProject         *projectusecases.CloseProject
	DeleteProject        *projectusecases.DeleteProject
	ExportProject        *projectusecases.ExportProject
	ImportProject        *projectusecases.ImportProject
	GetProjectToImport   *projectusecases.GetProjectToImport
	// Command Groups
	GetCommandGroups              *commandgroupusecases.GetCommandGroups
	CreateCommandGroup            *commandgroupusecases.CreateCommandGroup
	UpdateCommandGroup            *commandgroupusecases.UpdateCommandGroup
	DeleteCommandGroup            *commandgroupusecases.DeleteCommandGroup
	RemoveCommandFromCommandGroup *commandgroupusecases.RemoveCommandFromCommandGroup
	ReorderCommandGroups          *commandgroupusecases.ReorderCommandGroups
	RunCommandGroup               *commandgroupusecases.RunCommandGroup
	StopCommandGroup              *commandgroupusecases.StopCommandGroup
	// Commands
	GetCommands          *commandusecases.GetCommands
	AddCommand           *commandusecases.AddCommand
	DuplicateCommand     *commandusecases.DuplicateCommand
	RemoveCommand        *commandusecases.RemoveCommand
	EditCommand          *commandusecases.EditCommand
	ReorderCommands      *commandusecases.ReorderCommands
	RunCommand           *commandusecases.RunCommand
	StopCommand          *commandusecases.StopCommand
	GetRunningCommandIds *commandusecases.GetRunningCommandIds
}
