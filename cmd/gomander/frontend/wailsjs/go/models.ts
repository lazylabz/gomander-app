export namespace app {
	
	export interface CommandGroupJSONv1 {
	    name: string;
	    commandIds: string[];
	}
	export interface CommandJSONv1 {
	    id: string;
	    name: string;
	    command: string;
	    workingDirectory: string;
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
	}
	export interface ProjectExportJSONv1 {
	    version: number;
	    name: string;
	    commands: CommandJSONv1[];
	    commandGroups: CommandGroupJSONv1[];
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
	}
	export interface CommandGroup {
	    id: string;
	    projectId: string;
	    name: string;
	    commands: Command[];
	    position: number;
	}
	export interface EnvironmentPath {
	    id: string;
	    path: string;
	}
	export interface Config {
	    lastOpenedProjectId: string;
	    environmentPaths: EnvironmentPath[];
	}
	
	export interface Project {
	    id: string;
	    name: string;
	    workingDirectory: string;
	}

}

export namespace event {
	
	export enum Event {
	    GET_COMMANDS = "get_commands",
	    PROCESS_FINISHED = "process_finished",
	    PROCESS_STARTED = "process_started",
	    NEW_LOG_ENTRY = "new_log_entry",
	    ERROR_NOTIFICATION = "error_notification",
	    GET_USER_CONFIG = "get_user_config",
	    GET_COMMAND_GROUPS = "get_command_groups",
	}

}

