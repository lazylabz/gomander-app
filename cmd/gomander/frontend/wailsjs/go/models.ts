export namespace app {
	
	export interface UseCases {
	    GetUserConfig: any;
	    SaveUserConfig: any;
	    GetTranslation: any;
	    GetSupportedLanguages: any;
	    GetCurrentProject: any;
	    GetAvailableProjects: any;
	    OpenProject: any;
	    CreateProject: any;
	    EditProject: any;
	    CloseProject: any;
	    DeleteProject: any;
	    ExportProject: any;
	    ImportProject: any;
	    GetProjectToImport: any;
	    GetCommandGroups: any;
	    CreateCommandGroup: any;
	    UpdateCommandGroup: any;
	    DeleteCommandGroup: any;
	    RemoveCommandFromCommandGroup: any;
	    ReorderCommandGroups: any;
	    RunCommandGroup: any;
	    StopCommandGroup: any;
	    GetCommands: any;
	    AddCommand: any;
	    DuplicateCommand: any;
	    RemoveCommand: any;
	    EditCommand: any;
	    ReorderCommands: any;
	    RunCommand: any;
	    StopCommand: any;
	    GetRunningCommandIds: any;
	}
	export interface EventHandlers {
	    CleanCommandGroupsOnCommandDeleted: any;
	    CleanCommandGroupsOnProjectDeleted: any;
	    CleanCommandsOnProjectDeleted: any;
	    AddCommandToGroupOnCommandDuplicated: any;
	}
	export interface Dependencies {
	    Logger: any;
	    EventEmitter: any;
	    Runner: any;
	    CommandRepository: any;
	    CommandGroupRepository: any;
	    ProjectRepository: any;
	    ConfigRepository: any;
	    FsFacade: any;
	    RuntimeFacade: any;
	    EventBus: any;
	    EventHandlers: EventHandlers;
	    UseCases: UseCases;
	}
	

}

export namespace domain {
	
	export interface Command {
	    id: string;
	    projectId: string;
	    name: string;
	    command: string;
	    workingDirectory: string;
	    position: number;
	    link: string;
	    errorPatterns: string[];
	}
	export interface CommandGroup {
	    id: string;
	    projectId: string;
	    name: string;
	    commands: Command[];
	    position: number;
	}
	export interface CommandGroupJSONv1 {
	    id: string;
	    name: string;
	    commandIds: string[];
	}
	export interface CommandJSONv1 {
	    id: string;
	    name: string;
	    command: string;
	    workingDirectory: string;
	}
	export interface EnvironmentPath {
	    id: string;
	    path: string;
	}
	export interface Config {
	    lastOpenedProjectId: string;
	    environmentPaths: EnvironmentPath[];
	    logLineLimit: number;
	    locale: string;
	}
	
	export interface Localization {
	    "common.cancel": string;
	    "common.save": string;
	    "common.create": string;
	    "common.delete": string;
	    "common.edit": string;
	    "common.duplicate": string;
	    "common.continue": string;
	    "common.add": string;
	    "common.advanced": string;
	    "sidebar.project.close": string;
	    "sidebar.commands.title": string;
	    "sidebar.commands.add": string;
	    "sidebar.commands.removeFromGroup": string;
	    "sidebar.commandGroups.title": string;
	    "sidebar.commandGroups.add": string;
	    "sidebar.commandGroups.applyReorder": string;
	    "sidebar.commandGroups.startReorder": string;
	    "sidebar.create.title": string;
	    "sidebar.create.command": string;
	    "sidebar.create.commandGroup": string;
	    "sidebar.create.noCommandsTooltip": string;
	    "sidebar.version.current": string;
	    "sidebar.version.newAvailable": string;
	    "sidebar.version.latest": string;
	    "sidebar.version.checkError": string;
	    "projectSelection.openTitle": string;
	    "projectSelection.welcome": string;
	    "projectSelection.emptyState": string;
	    "projectSelection.createButton": string;
	    "projectSelection.importButton": string;
	    "projectSelection.moreOptions": string;
	    "projectSelection.importPackageJson": string;
	    "projectSelection.exportAction": string;
	    "logs.matches_one": string;
	    "logs.matches_other": string;
	    "settings.title": string;
	    "settings.tabs.user": string;
	    "settings.tabs.project": string;
	    "settings.saving.inProgress": string;
	    "settings.saving.done": string;
	    "aboutModal.version": string;
	    "aboutModal.newVersion": string;
	    "aboutModal.newVersionSubtitle": string;
	    "aboutModal.downloadUpdate": string;
	    "aboutModal.description": string;
	    "aboutModal.feedbackTitle": string;
	    "aboutModal.feedbackSubtitle": string;
	    "aboutModal.teamTitle": string;
	    "aboutModal.teamSubtitle": string;
	    "modal.createCommand.title": string;
	    "modal.editCommand.title": string;
	    "modal.createCommandGroup.title": string;
	    "modal.editCommandGroup.title": string;
	    "modal.createProject.title": string;
	    "modal.deleteProject.title": string;
	    "modal.deleteProject.description": string;
	    "modal.importProject.title": string;
	    "modal.importProject.description": string;
	    "modal.importProject.advancedTrigger": string;
	    "commandForm.nameLabel": string;
	    "commandForm.commandLabel": string;
	    "commandForm.errorPatternsLabel": string;
	    "commandForm.errorPatternsDescription": string;
	    "commandForm.linkLabel": string;
	    "commandForm.workingDirectoryLabel": string;
	    "commandForm.computedPath": string;
	    "commandForm.validation.nameRequired": string;
	    "commandForm.validation.commandRequired": string;
	    "commandGroupForm.nameLabel": string;
	    "commandGroupForm.commandsDescription": string;
	    "commandGroupForm.availableCommands": string;
	    "commandGroupForm.groupCommands": string;
	    "commandGroupForm.emptyAvailable": string;
	    "commandGroupForm.emptyGroup": string;
	    "commandGroupForm.validation.nameRequired": string;
	    "commandGroupForm.validation.commandsRequired": string;
	    "projectForm.nameLabel": string;
	    "projectForm.baseDirLabel": string;
	    "projectForm.commandsLabel": string;
	    "projectForm.commandGroupsLabel": string;
	    "projectForm.commandGroupsDisabledTooltip": string;
	    "projectForm.deletedCommand": string;
	    "projectForm.validation.nameRequired": string;
	    "projectForm.validation.baseDirRequired": string;
	    "userSettingsForm.envPathsTitle": string;
	    "userSettingsForm.envPathsDescription": string;
	    "userSettingsForm.envPathsHelpBody": string;
	    "userSettingsForm.envPathsHelpExample": string;
	    "userSettingsForm.preferencesTitle": string;
	    "userSettingsForm.preferencesDescription": string;
	    "userSettingsForm.languageLabel": string;
	    "userSettingsForm.languagePlaceholder": string;
	    "userSettingsForm.themeLabel": string;
	    "userSettingsForm.themePlaceholder": string;
	    "userSettingsForm.themeSystem": string;
	    "userSettingsForm.themeLight": string;
	    "userSettingsForm.themeDark": string;
	    "userSettingsForm.themeDescription": string;
	    "userSettingsForm.logLimitLabel": string;
	    "userSettingsForm.logLimitDescription": string;
	    "userSettingsForm.validation.pathEmpty": string;
	    "userSettingsForm.validation.logLimitMin": string;
	    "userSettingsForm.validation.logLimitMax": string;
	    "projectSettingsForm.sectionTitle": string;
	    "projectSettingsForm.sectionDescription": string;
	    "toast.command.runFailed": string;
	    "toast.command.stopFailed": string;
	    "toast.command.createSuccess": string;
	    "toast.command.createFailed": string;
	    "toast.command.updateSuccess": string;
	    "toast.command.updateFailed": string;
	    "toast.command.deleteSuccess": string;
	    "toast.command.deleteFailed": string;
	    "toast.command.duplicateSuccess": string;
	    "toast.command.duplicateFailed": string;
	    "toast.command.reorderFailed": string;
	    "toast.command.invalidDrag": string;
	    "toast.command.removeFromGroupSuccess": string;
	    "toast.command.removeFromGroupFailed": string;
	    "toast.commandGroup.createSuccess": string;
	    "toast.commandGroup.createFailed": string;
	    "toast.commandGroup.updateSuccess": string;
	    "toast.commandGroup.updateFailed": string;
	    "toast.commandGroup.deleteSuccess": string;
	    "toast.commandGroup.deleteFailed": string;
	    "toast.commandGroup.runFailed": string;
	    "toast.commandGroup.stopFailed": string;
	    "toast.commandGroup.reorderSuccess": string;
	    "toast.commandGroup.reorderFailed": string;
	    "toast.commandGroup.notFound": string;
	    "toast.commandGroup.cannotRemoveLast": string;
	    "toast.project.selectFailed": string;
	    "toast.project.importSuccess": string;
	    "toast.project.importFailed": string;
	    "toast.project.exportSuccess": string;
	    "toast.project.exportFailed": string;
	    "toast.project.openFolderAction": string;
	    "toast.settings.userSaveSuccess": string;
	    "toast.settings.userSaveFailed": string;
	    "toast.settings.projectSaveSuccess": string;
	    "toast.settings.projectSaveFailed": string;
	    "toast.version.checkError": string;
	}
	export interface Project {
	    id: string;
	    name: string;
	    workingDirectory: string;
	}
	export interface ProjectExportJSONv1 {
	    version: number;
	    name: string;
	    workingDirectory: string;
	    commands: CommandJSONv1[];
	    commandGroups: CommandGroupJSONv1[];
	}

}

export namespace event {
	
	export enum Event {
	    PROCESS_FINISHED = "process_finished",
	    PROCESS_STARTED = "process_started",
	    NEW_LOG_ENTRY = "new_log_entry",
	    COMMAND_GROUP_DELETED = "command_group_deleted",
	    COMMAND_ERROR_DETECTED = "command_error_detected",
	}

}

