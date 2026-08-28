import type {
	CommandExport,
	CommandGroupExport,
	ProjectExport,
} from "@/contracts/types.ts";

export class ProjectExportBuilder {
	private data: ProjectExport = {
		name: "Test Project",
		workingDirectory: "/app",
		commands: [],
		commandGroups: [],
	};

	withName(name: string): this {
		this.data.name = name;
		return this;
	}

	withWorkingDirectory(workingDirectory: string): this {
		this.data.workingDirectory = workingDirectory;
		return this;
	}

	withCommands(...commands: CommandExport[]): this {
		this.data.commands = commands;
		return this;
	}

	withCommandGroups(...commandGroups: CommandGroupExport[]): this {
		this.data.commandGroups = commandGroups;
		return this;
	}

	build(): ProjectExport {
		return {
			...this.data,
			commands: [...this.data.commands],
			commandGroups: [...this.data.commandGroups],
		};
	}
}
