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
	    Id: string;
	    ProjectId: string;
	    Name: string;
	    Commands: Command[];
	    Position: number;
	}
	export interface EnvironmentPath {
	    Id: string;
	    Path: string;
	}
	export interface Config {
	    LastOpenedProjectId: string;
	    EnvironmentPaths: EnvironmentPath[];
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
	    SUCCESS_NOTIFICATION = "success_notification",
	    GET_USER_CONFIG = "get_user_config",
	    GET_COMMAND_GROUPS = "get_command_groups",
	}

}

