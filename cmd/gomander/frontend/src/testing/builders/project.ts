import type { Project } from "@/contracts/types.ts";

export class ProjectBuilder {
	private data: Project = {
		id: crypto.randomUUID(),
		name: "Test Project",
		workingDirectory: "/app",
	};

	withId(id: string): this {
		this.data.id = id;
		return this;
	}

	withName(name: string): this {
		this.data.name = name;
		return this;
	}

	withWorkingDirectory(workingDirectory: string): this {
		this.data.workingDirectory = workingDirectory;
		return this;
	}

	build(): Project {
		return { ...this.data };
	}
}
