import type {
	Command,
	CommandGroup,
	Event,
	EventData,
	Localization,
	Project,
	ProjectExport,
	UserConfig,
} from "@/contracts/types.ts";

export type DataService = {
	addCommand: (command: Command) => Promise<void>;
	duplicateCommand: (commandId: string, groupId: string) => Promise<void>;
	editCommand: (command: Command) => Promise<void>;
	reorderCommands: (orderedCommandIds: string[]) => Promise<void>;
	getAvailableProjects: () => Promise<Project[]>;
	editCommandGroup: (commandGroup: CommandGroup) => Promise<void>;
	createCommandGroup: (commandGroup: CommandGroup) => Promise<void>;
	deleteCommandGroup: (groupId: string) => Promise<void>;
	reorderCommandGroups: (orderedGroupIds: string[]) => Promise<void>;
	removeCommandFromGroup: (commandId: string, groupId: string) => Promise<void>;
	runCommandGroup: (groupId: string) => Promise<void>;
	stopCommandGroup: (groupId: string) => Promise<void>;
	getCommandGroups: () => Promise<CommandGroup[]>;
	getCommands: () => Promise<Command[]>;
	getCurrentProject: () => Promise<Project | null>;
	getUserConfig: () => Promise<UserConfig>;
	removeCommand: (commandId: string) => Promise<void>;
	runCommand: (commandId: string) => Promise<void>;
	saveUserConfig: (userConfig: UserConfig) => Promise<void>;
	createProject: (project: Project) => Promise<void>;
	stopCommand: (commandId: string) => Promise<void>;
	openProject: (projectId: string) => Promise<void>;
	closeProject: () => Promise<void>;
	deleteProject: (projectId: string) => Promise<void>;
	exportProject: (projectId: string) => Promise<string>;
	importProject: (
		project: ProjectExport,
		name: string,
		workingDirectory: string,
	) => Promise<void>;
	getProjectToImport: () => Promise<ProjectExport>;
	getProjectToImportFromPackageJson: () => Promise<ProjectExport>;
	editProject: (project: Project) => Promise<void>;
};

export type HelpersService = {
	getComputedPath: (
		baseWorkingDirectory: string,
		workingDirectory: string,
	) => Promise<string>;
	isThereANewRelease: () => Promise<string>;
	getCurrentRelease: () => Promise<string>;
	downloadLatestRelease: (release: string) => Promise<string>;
	installLatestReleaseAndQuit: (binaryPath: string) => Promise<void>;
	getOs: () => Promise<string>;
	askForDirPath: () => Promise<string>;
	openFileFolder: (path: string) => Promise<void>;
};

export type EventService = {
	eventsOn: <E extends Event>(
		event: E,
		callback: (data: EventData[E]) => void,
	) => () => void;
	eventsOff: (event: Event, ...additionalEvents: Event[]) => void;
};

export type ExternalBrowserService = {
	browserOpenURL: (url: string) => void;
};

export type TranslationsService = {
	getTranslation: (language: string) => Promise<Localization>;
	getSupportedLanguages: () => Promise<string[]>;
};

export type BackendServices = {
	data: DataService;
	helpers: HelpersService;
	event: EventService;
	externalBrowser: ExternalBrowserService;
	translations: TranslationsService;
};
