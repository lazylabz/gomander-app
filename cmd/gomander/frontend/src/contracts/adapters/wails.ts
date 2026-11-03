import { Browser, Events } from "@wailsio/runtime";
import type { BackendServices } from "@/contracts/ports.ts";
import {
	AddCommandController,
	CheckForNewReleaseController,
	CloseProjectController,
	CreateCommandGroupController,
	CreateProjectController,
	DeleteCommandGroupController,
	DeleteProjectController,
	DownloadReleaseController,
	DuplicateCommandController,
	EditCommandController,
	EditProjectController,
	ExportProjectController,
	GetAvailableProjectsController,
	GetCommandGroupsController,
	GetCommandsController,
	GetCurrentProjectController,
	GetCurrentReleaseController,
	GetProjectToImportController,
	GetProjectToImportFromPackageJsonController,
	GetSupportedLanguagesController,
	GetTranslationController,
	GetUserConfigController,
	ImportProjectController,
	InstallReleaseAndQuitController,
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
} from "../../../bindings/gomander/cmd/gomander/wailscontrollers.ts";
import {
	AskForDirPath,
	OpenFileFolder,
} from "../../../bindings/gomander/internal/uihelpers/fs/uifshelper.ts";
import { GetOs } from "../../../bindings/gomander/internal/uihelpers/os_internal/uioshelper.ts";
import { GetComputedPath } from "../../../bindings/gomander/internal/uihelpers/path/uipathhelper.ts";
import type { EventData, WithoutNullSlices } from "../types.ts";

// v3 types a call answering with a Go slice or pointer as `| null`; the contract
// answers the way v2 typed it, so the call is only recast, never wrapped.
const asContract = <Args extends unknown[], Result>(
	call: (...args: Args) => PromiseLike<Result>,
) =>
	call as unknown as (
		...args: Args
	) => Promise<WithoutNullSlices<NonNullable<Result>>>;

export const wailsBackend = {
	data: {
		addCommand: AddCommandController,
		duplicateCommand: DuplicateCommandController,
		editCommand: EditCommandController,
		reorderCommands: ReorderCommandsController,
		getAvailableProjects: asContract(GetAvailableProjectsController),
		editCommandGroup: UpdateCommandGroupController,
		createCommandGroup: CreateCommandGroupController,
		deleteCommandGroup: DeleteCommandGroupController,
		reorderCommandGroups: ReorderCommandGroupsController,
		removeCommandFromGroup: RemoveCommandFromCommandGroupController,
		runCommandGroup: RunCommandGroupController,
		stopCommandGroup: StopCommandGroupController,
		getCommandGroups: asContract(GetCommandGroupsController),
		getCommands: asContract(GetCommandsController),
		getCurrentProject: GetCurrentProjectController,
		getUserConfig: asContract(GetUserConfigController),
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
		getProjectToImport: asContract(GetProjectToImportController),
		getProjectToImportFromPackageJson: asContract(
			GetProjectToImportFromPackageJsonController,
		),
		editProject: EditProjectController,
		getCurrentRelease: GetCurrentReleaseController,
		checkForNewRelease: CheckForNewReleaseController,
		downloadRelease: DownloadReleaseController,
		installReleaseAndQuit: InstallReleaseAndQuitController,
	},
	helpers: {
		getComputedPath: GetComputedPath,
		getOs: GetOs,
		askForDirPath: AskForDirPath,
		openFileFolder: OpenFileFolder,
	},
	event: {
		eventsOn: (event, callback) =>
			Events.On(event, ({ data }) => callback(data as EventData[typeof event])),
		eventsOff: (event, ...additionalEvents) =>
			Events.Off(event, ...additionalEvents),
	},
	externalBrowser: {
		browserOpenURL: (url) => {
			void Browser.OpenURL(url);
		},
	},
	translations: {
		getTranslation: asContract(GetTranslationController),
		getSupportedLanguages: asContract(GetSupportedLanguagesController),
	},
} satisfies BackendServices;
