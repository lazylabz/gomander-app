import type {
	CommandBlueprint,
	CommandGroupBlueprint,
	ProjectBlueprint,
} from "@/contracts/types.ts";

export class ProjectBlueprintBuilder {
	private data: ProjectBlueprint = {
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

	withCommands(...commands: CommandBlueprint[]): this {
		this.data.commands = commands;
		return this;
	}

	withCommandGroups(...commandGroups: CommandGroupBlueprint[]): this {
		this.data.commandGroups = commandGroups;
		return this;
	}

	build(): ProjectBlueprint {
		return {
			...this.data,
			commands: [...this.data.commands],
			commandGroups: [...this.data.commandGroups],
		};
	}
}
