export namespace config {
	
	export class UserConfig {
	    environmentPaths: string[];
	    lastOpenedProjectId: string;
	
	    static createFrom(source: any = {}) {
	        return new UserConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.environmentPaths = source["environmentPaths"];
	        this.lastOpenedProjectId = source["lastOpenedProjectId"];
	    }
	}

}

export namespace event {
	
	export enum Event {
	    GET_COMMANDS = "get_commands",
	    PROCESS_FINISHED = "process_finished",
	    NEW_LOG_ENTRY = "new_log_entry",
	    ERROR_NOTIFICATION = "error_notification",
	    SUCCESS_NOTIFICATION = "success_notification",
	    GET_USER_CONFIG = "get_user_config",
	    GET_COMMAND_GROUPS = "get_command_groups",
	}

}

export namespace project {
	
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
	export class CommandGroup {
	    id: string;
	    name: string;
	    commands: string[];
	
	    static createFrom(source: any = {}) {
	        return new CommandGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.commands = source["commands"];
	    }
	}
	export class Project {
	    id: string;
	    name: string;
	    commands: Record<string, Command>;
	    commandGroups: CommandGroup[];
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.commands = this.convertValues(source["commands"], Command, true);
	        this.commandGroups = this.convertValues(source["commandGroups"], CommandGroup);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

