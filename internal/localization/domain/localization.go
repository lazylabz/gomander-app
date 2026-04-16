package domain

type Localization struct {
	// common
	CommonCancel    string `json:"common.cancel"`
	CommonSave      string `json:"common.save"`
	CommonCreate    string `json:"common.create"`
	CommonDelete    string `json:"common.delete"`
	CommonEdit      string `json:"common.edit"`
	CommonDuplicate string `json:"common.duplicate"`
	CommonContinue  string `json:"common.continue"`
	CommonAdd       string `json:"common.add"`
	CommonAdvanced  string `json:"common.advanced"`

	// sidebar
	SidebarProjectClose                string `json:"sidebar.project.close"`
	SidebarCommandsTitle               string `json:"sidebar.commands.title"`
	SidebarCommandsAdd                 string `json:"sidebar.commands.add"`
	SidebarCommandsRemoveFromGroup     string `json:"sidebar.commands.removeFromGroup"`
	SidebarCommandGroupsTitle          string `json:"sidebar.commandGroups.title"`
	SidebarCommandGroupsAdd            string `json:"sidebar.commandGroups.add"`
	SidebarCommandGroupsApplyReorder   string `json:"sidebar.commandGroups.applyReorder"`
	SidebarCommandGroupsStartReorder   string `json:"sidebar.commandGroups.startReorder"`
	SidebarCreateTitle                 string `json:"sidebar.create.title"`
	SidebarCreateCommand               string `json:"sidebar.create.command"`
	SidebarCreateCommandGroup          string `json:"sidebar.create.commandGroup"`
	SidebarCreateNoCommandsTooltip     string `json:"sidebar.create.noCommandsTooltip"`
	SidebarVersionCurrent              string `json:"sidebar.version.current"`
	SidebarVersionNewAvailable         string `json:"sidebar.version.newAvailable"`
	SidebarVersionLatest               string `json:"sidebar.version.latest"`
	SidebarVersionCheckError           string `json:"sidebar.version.checkError"`

	// projectSelection
	ProjectSelectionOpenTitle        string `json:"projectSelection.openTitle"`
	ProjectSelectionWelcome          string `json:"projectSelection.welcome"`
	ProjectSelectionEmptyState       string `json:"projectSelection.emptyState"`
	ProjectSelectionCreateButton     string `json:"projectSelection.createButton"`
	ProjectSelectionImportButton     string `json:"projectSelection.importButton"`
	ProjectSelectionMoreOptions      string `json:"projectSelection.moreOptions"`
	ProjectSelectionImportPackageJson string `json:"projectSelection.importPackageJson"`
	ProjectSelectionExportAction     string `json:"projectSelection.exportAction"`

	// logs
	LogsMatchesOne   string `json:"logs.matches_one"`
	LogsMatchesOther string `json:"logs.matches_other"`

	// settings screen
	SettingsTitle           string `json:"settings.title"`
	SettingsTabsUser        string `json:"settings.tabs.user"`
	SettingsTabsProject     string `json:"settings.tabs.project"`
	SettingsSavingInProgress string `json:"settings.saving.inProgress"`
	SettingsSavingDone      string `json:"settings.saving.done"`

	// aboutModal
	AboutModalVersion            string `json:"aboutModal.version"`
	AboutModalNewVersion         string `json:"aboutModal.newVersion"`
	AboutModalNewVersionSubtitle string `json:"aboutModal.newVersionSubtitle"`
	AboutModalDownloadUpdate     string `json:"aboutModal.downloadUpdate"`
	AboutModalDescription        string `json:"aboutModal.description"`
	AboutModalFeedbackTitle      string `json:"aboutModal.feedbackTitle"`
	AboutModalFeedbackSubtitle   string `json:"aboutModal.feedbackSubtitle"`
	AboutModalTeamTitle          string `json:"aboutModal.teamTitle"`
	AboutModalTeamSubtitle       string `json:"aboutModal.teamSubtitle"`

	// modal chrome
	ModalCreateCommandTitle         string `json:"modal.createCommand.title"`
	ModalEditCommandTitle           string `json:"modal.editCommand.title"`
	ModalCreateCommandGroupTitle    string `json:"modal.createCommandGroup.title"`
	ModalEditCommandGroupTitle      string `json:"modal.editCommandGroup.title"`
	ModalCreateProjectTitle         string `json:"modal.createProject.title"`
	ModalDeleteProjectTitle         string `json:"modal.deleteProject.title"`
	ModalDeleteProjectDescription   string `json:"modal.deleteProject.description"`
	ModalImportProjectTitle         string `json:"modal.importProject.title"`
	ModalImportProjectDescription   string `json:"modal.importProject.description"`
	ModalImportProjectAdvancedTrigger string `json:"modal.importProject.advancedTrigger"`

	// commandForm
	CommandFormNameLabel                string `json:"commandForm.nameLabel"`
	CommandFormCommandLabel             string `json:"commandForm.commandLabel"`
	CommandFormErrorPatternsLabel       string `json:"commandForm.errorPatternsLabel"`
	CommandFormErrorPatternsDescription string `json:"commandForm.errorPatternsDescription"`
	CommandFormLinkLabel                string `json:"commandForm.linkLabel"`
	CommandFormWorkingDirectoryLabel    string `json:"commandForm.workingDirectoryLabel"`
	CommandFormComputedPath             string `json:"commandForm.computedPath"`
	CommandFormValidationNameRequired   string `json:"commandForm.validation.nameRequired"`
	CommandFormValidationCommandRequired string `json:"commandForm.validation.commandRequired"`

	// commandGroupForm
	CommandGroupFormNameLabel                  string `json:"commandGroupForm.nameLabel"`
	CommandGroupFormCommandsDescription        string `json:"commandGroupForm.commandsDescription"`
	CommandGroupFormAvailableCommands          string `json:"commandGroupForm.availableCommands"`
	CommandGroupFormGroupCommands              string `json:"commandGroupForm.groupCommands"`
	CommandGroupFormEmptyAvailable             string `json:"commandGroupForm.emptyAvailable"`
	CommandGroupFormEmptyGroup                 string `json:"commandGroupForm.emptyGroup"`
	CommandGroupFormValidationNameRequired     string `json:"commandGroupForm.validation.nameRequired"`
	CommandGroupFormValidationCommandsRequired string `json:"commandGroupForm.validation.commandsRequired"`

	// projectForm
	ProjectFormNameLabel                  string `json:"projectForm.nameLabel"`
	ProjectFormBaseDirLabel               string `json:"projectForm.baseDirLabel"`
	ProjectFormCommandsLabel              string `json:"projectForm.commandsLabel"`
	ProjectFormCommandGroupsLabel         string `json:"projectForm.commandGroupsLabel"`
	ProjectFormCommandGroupsDisabledTooltip string `json:"projectForm.commandGroupsDisabledTooltip"`
	ProjectFormDeletedCommand             string `json:"projectForm.deletedCommand"`
	ProjectFormValidationNameRequired     string `json:"projectForm.validation.nameRequired"`
	ProjectFormValidationBaseDirRequired  string `json:"projectForm.validation.baseDirRequired"`

	// userSettingsForm
	UserSettingsFormEnvPathsTitle          string `json:"userSettingsForm.envPathsTitle"`
	UserSettingsFormEnvPathsDescription    string `json:"userSettingsForm.envPathsDescription"`
	UserSettingsFormEnvPathsHelpBody       string `json:"userSettingsForm.envPathsHelpBody"`
	UserSettingsFormEnvPathsHelpExample    string `json:"userSettingsForm.envPathsHelpExample"`
	UserSettingsFormPreferencesTitle       string `json:"userSettingsForm.preferencesTitle"`
	UserSettingsFormPreferencesDescription string `json:"userSettingsForm.preferencesDescription"`
	UserSettingsFormLanguageLabel          string `json:"userSettingsForm.languageLabel"`
	UserSettingsFormLanguagePlaceholder    string `json:"userSettingsForm.languagePlaceholder"`
	UserSettingsFormThemeLabel             string `json:"userSettingsForm.themeLabel"`
	UserSettingsFormThemePlaceholder       string `json:"userSettingsForm.themePlaceholder"`
	UserSettingsFormThemeSystem            string `json:"userSettingsForm.themeSystem"`
	UserSettingsFormThemeLight             string `json:"userSettingsForm.themeLight"`
	UserSettingsFormThemeDark              string `json:"userSettingsForm.themeDark"`
	UserSettingsFormThemeDescription       string `json:"userSettingsForm.themeDescription"`
	UserSettingsFormLogLimitLabel          string `json:"userSettingsForm.logLimitLabel"`
	UserSettingsFormLogLimitDescription    string `json:"userSettingsForm.logLimitDescription"`
	UserSettingsFormValidationPathEmpty    string `json:"userSettingsForm.validation.pathEmpty"`
	UserSettingsFormValidationLogLimitMin  string `json:"userSettingsForm.validation.logLimitMin"`
	UserSettingsFormValidationLogLimitMax  string `json:"userSettingsForm.validation.logLimitMax"`

	// projectSettingsForm
	ProjectSettingsFormSectionTitle       string `json:"projectSettingsForm.sectionTitle"`
	ProjectSettingsFormSectionDescription string `json:"projectSettingsForm.sectionDescription"`

	// toast.command
	ToastCommandRunFailed             string `json:"toast.command.runFailed"`
	ToastCommandStopFailed            string `json:"toast.command.stopFailed"`
	ToastCommandCreateSuccess         string `json:"toast.command.createSuccess"`
	ToastCommandCreateFailed          string `json:"toast.command.createFailed"`
	ToastCommandUpdateSuccess         string `json:"toast.command.updateSuccess"`
	ToastCommandUpdateFailed          string `json:"toast.command.updateFailed"`
	ToastCommandDeleteSuccess         string `json:"toast.command.deleteSuccess"`
	ToastCommandDeleteFailed          string `json:"toast.command.deleteFailed"`
	ToastCommandDuplicateSuccess      string `json:"toast.command.duplicateSuccess"`
	ToastCommandDuplicateFailed       string `json:"toast.command.duplicateFailed"`
	ToastCommandReorderFailed         string `json:"toast.command.reorderFailed"`
	ToastCommandInvalidDrag           string `json:"toast.command.invalidDrag"`
	ToastCommandRemoveFromGroupSuccess string `json:"toast.command.removeFromGroupSuccess"`
	ToastCommandRemoveFromGroupFailed  string `json:"toast.command.removeFromGroupFailed"`

	// toast.commandGroup
	ToastCommandGroupCreateSuccess     string `json:"toast.commandGroup.createSuccess"`
	ToastCommandGroupCreateFailed      string `json:"toast.commandGroup.createFailed"`
	ToastCommandGroupUpdateSuccess     string `json:"toast.commandGroup.updateSuccess"`
	ToastCommandGroupUpdateFailed      string `json:"toast.commandGroup.updateFailed"`
	ToastCommandGroupDeleteSuccess     string `json:"toast.commandGroup.deleteSuccess"`
	ToastCommandGroupDeleteFailed      string `json:"toast.commandGroup.deleteFailed"`
	ToastCommandGroupRunFailed         string `json:"toast.commandGroup.runFailed"`
	ToastCommandGroupStopFailed        string `json:"toast.commandGroup.stopFailed"`
	ToastCommandGroupReorderSuccess    string `json:"toast.commandGroup.reorderSuccess"`
	ToastCommandGroupReorderFailed     string `json:"toast.commandGroup.reorderFailed"`
	ToastCommandGroupNotFound          string `json:"toast.commandGroup.notFound"`
	ToastCommandGroupCannotRemoveLast  string `json:"toast.commandGroup.cannotRemoveLast"`

	// toast.project
	ToastProjectSelectFailed     string `json:"toast.project.selectFailed"`
	ToastProjectImportSuccess    string `json:"toast.project.importSuccess"`
	ToastProjectImportFailed     string `json:"toast.project.importFailed"`
	ToastProjectExportSuccess    string `json:"toast.project.exportSuccess"`
	ToastProjectExportFailed     string `json:"toast.project.exportFailed"`
	ToastProjectOpenFolderAction string `json:"toast.project.openFolderAction"`

	// toast.settings
	ToastSettingsUserSaveSuccess    string `json:"toast.settings.userSaveSuccess"`
	ToastSettingsUserSaveFailed     string `json:"toast.settings.userSaveFailed"`
	ToastSettingsProjectSaveSuccess string `json:"toast.settings.projectSaveSuccess"`
	ToastSettingsProjectSaveFailed  string `json:"toast.settings.projectSaveFailed"`

	// toast.version
	ToastVersionCheckError string `json:"toast.version.checkError"`
}
