export namespace main {
	
	export enum Event {
	    GET_COMMANDS = "get_commands",
	    PROCESS_FINISHED = "process_finished",
	    NEW_LOG_ENTRY = "new_log_entry",
	    ERROR_NOTIFICATION = "error_notification",
	    SUCCESS_NOTIFICATION = "success_notification",
	    GET_USER_CONFIG = "get_user_config",
	}
	export class Command {
	    id: string;
	    name: string;
	    command: string;
	    workingDirectory: string;
	
	    static createFrom(source: any = {}) {
	        return new Command(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.command = source["command"];
	        this.workingDirectory = source["workingDirectory"];
	    }
	}
	export class UserConfig {
	    extraPaths: string[];
	
	    static createFrom(source: any = {}) {
	        return new UserConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.extraPaths = source["extraPaths"];
	    }
	}

}

