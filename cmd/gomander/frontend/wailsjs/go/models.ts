export namespace config {
	
	export interface UserConfig {
	    environmentPaths: string[];
	    lastOpenedProjectId: string;
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

export namespace project {
	
	export interface Command {
	    id: string;
	    name: string;
	    command: string;
	    workingDirectory: string;
	}
	export interface CommandGroup {
	    id: string;
	    name: string;
	    commands: string[];
	}
	export interface ExportableProject {
	    id: string;
	    name: string;
	    commands: Record<string, Command>;
	    commandGroups: CommandGroup[];
	}
	export interface Project {
	    id: string;
	    name: string;
	    baseWorkingDirectory: string;
	    commands: Record<string, Command>;
	    commandGroups: CommandGroup[];
	}

}

