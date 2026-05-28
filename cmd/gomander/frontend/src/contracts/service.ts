import { AskForDirPath, OpenFileFolder } from "../../wailsjs/go/fs/UIFsHelper";
import {
	AddCommandController,
	CloseProjectController,
	CreateCommandGroupController,
	CreateProjectController,
	DeleteCommandGroupController,
	DeleteProjectController,
	DuplicateCommandController,
	EditCommandController,
	EditProjectController,
	ExportProjectController,
	GetAvailableProjectsController,
	GetCommandGroupsController,
	GetCommandsController,
	GetCurrentProjectController,
	GetProjectToImportController,
	GetProjectToImportFromPackageJsonController,
	GetSupportedLanguagesController,
	GetTranslationController,
	GetUserConfigController,
	ImportProjectController,
	OpenProjectController,
	RemoveCommandController,
	RemoveCommandFromCommandGroupController,
	ReorderCommandGroupsController,
	ReorderCommandsController,
	RunCommandController,
	RunCommandGroupController,
	SaveUserConfigController,
	StopCommandController,
	StopCommandGroupController,
	UpdateCommandGroupController,
} from "../../wailsjs/go/main/WailsControllers";
import type { domain } from "../../wailsjs/go/models.ts";
import { GetOs } from "../../wailsjs/go/os_internal/UIOsHelper";
import { GetComputedPath } from "../../wailsjs/go/path/UiPathHelper";
import {
	DownloadLatestRelease,
	GetCurrentRelease,
	InstallLatestReleaseAndQuit,
	IsThereANewRelease,
} from "../../wailsjs/go/releases/ReleaseHelper";
import { BrowserOpenURL, EventsOff, EventsOn } from "../../wailsjs/runtime";

export const dataService = {
	addCommand: AddCommandController,
	duplicateCommand: DuplicateCommandController,
	editCommand: EditCommandController,
	reorderCommands: ReorderCommandsController,
	getAvailableProjects: GetAvailableProjectsController,
	editCommandGroup: UpdateCommandGroupController,
	createCommandGroup: CreateCommandGroupController,
	deleteCommandGroup: DeleteCommandGroupController,
	reorderCommandGroups: ReorderCommandGroupsController,
	removeCommandFromGroup: RemoveCommandFromCommandGroupController,
	runCommandGroup: RunCommandGroupController,
	stopCommandGroup: StopCommandGroupController,
	getCommandGroups: GetCommandGroupsController,
	getCommands: GetCommandsController,
	getCurrentProject:
		GetCurrentProjectController as () => Promise<domain.Project | null>,
	getUserConfig: GetUserConfigController,
	removeCommand: RemoveCommandController,
	runCommand: RunCommandController,
	saveUserConfig: SaveUserConfigController,
	createProject: CreateProjectController,
	stopCommand: StopCommandController,
	openProject: OpenProjectController,
	closeProject: CloseProjectController,
	deleteProject: DeleteProjectController,
	exportProject: ExportProjectController,
	importProject: ImportProjectController,
	getProjectToImport: GetProjectToImportController,
	getProjectToImportFromPackageJson:
		GetProjectToImportFromPackageJsonController,
	editProject: EditProjectController,
};

export const helpersService = {
	getComputedPath: GetComputedPath,
	isThereANewRelease: IsThereANewRelease,
	getCurrentRelease: GetCurrentRelease,
	downloadLatestRelease: DownloadLatestRelease,
	installLatestReleaseAndQuit: InstallLatestReleaseAndQuit,
	getOs: GetOs,
	askForDirPath: AskForDirPath,
	openFileFolder: OpenFileFolder,
};

export const eventService = {
	eventsOn: EventsOn,
	eventsOff: EventsOff,
};

export const externalBrowserService = {
	browserOpenURL: BrowserOpenURL,
};

export const translationsService = {
	getTranslation: GetTranslationController,
	getSupportedLanguages: GetSupportedLanguagesController,
};
