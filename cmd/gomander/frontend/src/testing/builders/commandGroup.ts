import type { Command, CommandGroup } from "@/contracts/types.ts";

export class CommandGroupBuilder {
	private data: CommandGroup = {
		id: crypto.randomUUID(),
		projectId: crypto.randomUUID(),
		name: "Test Command Group",
		position: 0,
		commandIds: [],
	};

	withId(id: string): this {
		this.data.id = id;
		return this;
	}

	withProjectId(projectId: string): this {
		this.data.projectId = projectId;
		return this;
	}

	withName(name: string): this {
		this.data.name = name;
		return this;
	}

	withPosition(position: number): this {
		this.data.position = position;
		return this;
	}

	withCommands(...commands: Command[]): this {
		this.data.commandIds = commands.map((command) => command.id);
		return this;
	}

	build(): CommandGroup {
		return { ...this.data, commandIds: [...this.data.commandIds] };
	}
}
