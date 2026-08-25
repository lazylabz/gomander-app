import type { Command } from "@/contracts/types.ts";

export class CommandBuilder {
	private data: Command = {
		id: crypto.randomUUID(),
		projectId: crypto.randomUUID(),
		name: "Test Command",
		command: "echo 'hello'",
		workingDirectory: "/app",
		position: 0,
		link: "",
		errorPatterns: [],
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

	withCommand(command: string): this {
		this.data.command = command;
		return this;
	}

	withWorkingDirectory(workingDirectory: string): this {
		this.data.workingDirectory = workingDirectory;
		return this;
	}

	withPosition(position: number): this {
		this.data.position = position;
		return this;
	}

	withLink(link: string): this {
		this.data.link = link;
		return this;
	}

	withErrorPatterns(...errorPatterns: string[]): this {
		this.data.errorPatterns = errorPatterns;
		return this;
	}

	build(): Command {
		return { ...this.data, errorPatterns: [...this.data.errorPatterns] };
	}
}
