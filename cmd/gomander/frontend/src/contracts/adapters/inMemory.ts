// The values this fake derives (duplicate ids, export paths, computed paths) are
// its own convention, not the backend's: assert on what a test put in, never on
// the shape of what the fake made up.
import type { BackendServices } from "@/contracts/ports.ts";
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

export type InMemoryBackendState = {
	commands: Command[];
	commandGroups: CommandGroup[];
	projects: Project[];
	currentProject: Project | null;
	userConfig: UserConfig;
	runningCommandIds: string[];
	runningGroupIds: string[];
	projectToImport: ProjectExport;
	importedProjects: ProjectExport[];
	translations: Record<string, Localization>;
	supportedLanguages: string[];
	os: string;
	currentRelease: string;
	newRelease: string;
	pickedDirPath: string;
	downloadedReleases: string[];
	installedBinaryPath: string | null;
	openedUrls: string[];
	openedFolders: string[];
	exportedProjectIds: string[];
};

export type InMemoryBackend = BackendServices & {
	state: InMemoryBackendState;
	emit: <E extends Event>(event: E, data: EventData[E]) => void;
};

const emptyProjectExport: ProjectExport = {
	name: "",
	workingDirectory: "",
	commands: [],
	commandGroups: [],
};

const createState = (): InMemoryBackendState => ({
	commands: [],
	commandGroups: [],
	projects: [],
	currentProject: null,
	userConfig: {
		lastOpenedProjectId: "",
		environmentPaths: [],
		locale: "en",
	},
	runningCommandIds: [],
	runningGroupIds: [],
	projectToImport: emptyProjectExport,
	importedProjects: [],
	translations: {},
	supportedLanguages: ["en"],
	os: "darwin",
	currentRelease: "v1.0.0",
	newRelease: "",
	pickedDirPath: "/picked/dir",
	downloadedReleases: [],
	installedBinaryPath: null,
	openedUrls: [],
	openedFolders: [],
	exportedProjectIds: [],
});

export const createInMemoryBackend = (
	initialState: Partial<InMemoryBackendState> = {},
): InMemoryBackend => {
	// Cloned on the way in and on the way out: the real bindings deserialize fresh
	// objects, so neither the caller's fixtures nor the fake's state can be
	// corrupted by whoever holds the other side.
	const state: InMemoryBackendState = structuredClone({
		...createState(),
		...initialState,
	});
	const snapshot = <T>(value: T): T => structuredClone(value);

	// Payload types are per-event, so the map cannot be typed without collapsing
	// them; the cast in `eventsOn` is where the per-event type is given up.
	const listeners = new Map<Event, ((data: unknown) => void)[]>();

	const reorderBy = <T extends { id: string }>(
		items: T[],
		orderedIds: string[],
	): T[] => {
		// The backend is given the full order; a partial list would silently sort
		// the unlisted items to the front instead of failing.
		const unlisted = items.filter((item) => !orderedIds.includes(item.id));
		if (unlisted.length > 0) {
			throw new Error(
				`Reorder left out ${unlisted.map((item) => item.id).join(", ")}`,
			);
		}

		return [...items].sort(
			(a, b) => orderedIds.indexOf(a.id) - orderedIds.indexOf(b.id),
		);
	};

	return {
		state,

		emit: (event, data) => {
			for (const listener of listeners.get(event) ?? []) {
				listener(data);
			}
		},

		data: {
			addCommand: async (command) => {
				state.commands.push(snapshot(command));
			},
			duplicateCommand: async (commandId, groupId) => {
				const command = state.commands.find((c) => c.id === commandId);
				if (!command) {
					throw new Error(`Command ${commandId} not found`);
				}

				const duplicate = { ...command, id: `${command.id}-copy` };
				state.commands.push(duplicate);

				const group = state.commandGroups.find((g) => g.id === groupId);
				group?.commands.push(duplicate);
			},
			editCommand: async (command) => {
				state.commands = state.commands.map((c) =>
					c.id === command.id ? snapshot(command) : c,
				);
				// A group carries a copy of every command it holds, so an edit shows
				// up there too.
				state.commandGroups = state.commandGroups.map((g) => ({
					...g,
					commands: g.commands.map((c) =>
						c.id === command.id ? snapshot(command) : c,
					),
				}));
			},
			reorderCommands: async (orderedCommandIds) => {
				state.commands = reorderBy(state.commands, orderedCommandIds);
			},
			getAvailableProjects: async () => snapshot(state.projects),
			editCommandGroup: async (commandGroup) => {
				state.commandGroups = state.commandGroups.map((g) =>
					g.id === commandGroup.id ? snapshot(commandGroup) : g,
				);
			},
			createCommandGroup: async (commandGroup) => {
				state.commandGroups.push(snapshot(commandGroup));
			},
			deleteCommandGroup: async (groupId) => {
				state.commandGroups = state.commandGroups.filter(
					(g) => g.id !== groupId,
				);
			},
			reorderCommandGroups: async (orderedGroupIds) => {
				state.commandGroups = reorderBy(state.commandGroups, orderedGroupIds);
			},
			removeCommandFromGroup: async (commandId, groupId) => {
				state.commandGroups = state.commandGroups.map((g) =>
					g.id === groupId
						? { ...g, commands: g.commands.filter((c) => c.id !== commandId) }
						: g,
				);
			},
			runCommandGroup: async (groupId) => {
				state.runningGroupIds.push(groupId);
			},
			stopCommandGroup: async (groupId) => {
				state.runningGroupIds = state.runningGroupIds.filter(
					(id) => id !== groupId,
				);
			},
			getCommandGroups: async () => snapshot(state.commandGroups),
			getCommands: async () => snapshot(state.commands),
			getCurrentProject: async () => snapshot(state.currentProject),
			getUserConfig: async () => snapshot(state.userConfig),
			removeCommand: async (commandId) => {
				state.commands = state.commands.filter((c) => c.id !== commandId);
				// Mirrors the backend event handler that prunes deleted commands from
				// every group they belonged to.
				state.commandGroups = state.commandGroups.map((g) => ({
					...g,
					commands: g.commands.filter((c) => c.id !== commandId),
				}));
			},
			runCommand: async (commandId) => {
				state.runningCommandIds.push(commandId);
			},
			saveUserConfig: async (userConfig) => {
				state.userConfig = snapshot(userConfig);
			},
			createProject: async (project) => {
				state.projects.push(snapshot(project));
			},
			stopCommand: async (commandId) => {
				state.runningCommandIds = state.runningCommandIds.filter(
					(id) => id !== commandId,
				);
			},
			openProject: async (projectId) => {
				const project = state.projects.find((p) => p.id === projectId);
				if (!project) {
					throw new Error(`Project ${projectId} not found`);
				}
				state.currentProject = project;
			},
			closeProject: async () => {
				state.currentProject = null;
			},
			deleteProject: async (projectId) => {
				state.projects = state.projects.filter((p) => p.id !== projectId);
				if (state.currentProject?.id === projectId) {
					state.currentProject = null;
				}
			},
			exportProject: async (projectId) => {
				state.exportedProjectIds.push(projectId);
				return `/exports/${projectId}.json`;
			},
			importProject: async (project, name, workingDirectory) => {
				state.importedProjects.push(snapshot(project));
				state.projects.push({
					id: `imported-${name}`,
					name,
					workingDirectory,
				});
			},
			getProjectToImport: async () => snapshot(state.projectToImport),
			getProjectToImportFromPackageJson: async () =>
				snapshot(state.projectToImport),
			editProject: async (project) => {
				state.projects = state.projects.map((p) =>
					p.id === project.id ? snapshot(project) : p,
				);
				if (state.currentProject?.id === project.id) {
					state.currentProject = snapshot(project);
				}
			},
		},

		helpers: {
			getComputedPath: async (baseWorkingDirectory, workingDirectory) =>
				[baseWorkingDirectory, workingDirectory].filter(Boolean).join("/"),
			isThereANewRelease: async () => state.newRelease,
			getCurrentRelease: async () => state.currentRelease,
			downloadLatestRelease: async (release) => {
				state.downloadedReleases.push(release);
				return `/downloads/${release}`;
			},
			installLatestReleaseAndQuit: async (binaryPath) => {
				state.installedBinaryPath = binaryPath;
			},
			getOs: async () => state.os,
			askForDirPath: async () => state.pickedDirPath,
			openFileFolder: async (path) => {
				state.openedFolders.push(path);
			},
		},

		event: {
			eventsOn: (event, callback) => {
				const listener = callback as (data: unknown) => void;
				listeners.set(event, [...(listeners.get(event) ?? []), listener]);

				return () => {
					listeners.set(
						event,
						(listeners.get(event) ?? []).filter((l) => l !== listener),
					);
				};
			},
			eventsOff: (event, ...additionalEvents) => {
				for (const eventName of [event, ...additionalEvents]) {
					listeners.delete(eventName);
				}
			},
		},

		externalBrowser: {
			browserOpenURL: (url) => {
				state.openedUrls.push(url);
			},
		},

		translations: {
			getTranslation: async (language) => {
				const translation = state.translations[language];
				if (!translation) {
					throw new Error(`No translation registered for "${language}"`);
				}
				return snapshot(translation);
			},
			getSupportedLanguages: async () => snapshot(state.supportedLanguages),
		},
	};
};
