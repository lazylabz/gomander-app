export namespace app {
	
	export interface UseCases {
	    GetUserConfig: any;
	    SaveUserConfig: any;
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

