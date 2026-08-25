import type { Command } from "@/contracts/types.ts";

export type CommandBuilder = {
	withId: (id: string) => CommandBuilder;
	withProjectId: (projectId: string) => CommandBuilder;
	withName: (name: string) => CommandBuilder;
	withCommand: (command: string) => CommandBuilder;
	withWorkingDirectory: (workingDirectory: string) => CommandBuilder;
	withPosition: (position: number) => CommandBuilder;
	withLink: (link: string) => CommandBuilder;
	withErrorPatterns: (...errorPatterns: string[]) => CommandBuilder;
	build: () => Command;
};

const builder = (data: Command): CommandBuilder => ({
	withId: (id) => builder({ ...data, id }),
	withProjectId: (projectId) => builder({ ...data, projectId }),
	withName: (name) => builder({ ...data, name }),
	withCommand: (command) => builder({ ...data, command }),
	withWorkingDirectory: (workingDirectory) =>
		builder({ ...data, workingDirectory }),
	withPosition: (position) => builder({ ...data, position }),
	withLink: (link) => builder({ ...data, link }),
	withErrorPatterns: (...errorPatterns) => builder({ ...data, errorPatterns }),
	build: () => ({ ...data }),
});

export const newCommandBuilder = (): CommandBuilder =>
	builder({
		id: crypto.randomUUID(),
		projectId: crypto.randomUUID(),
		name: "Test Command",
		command: "echo 'hello'",
		workingDirectory: "/app",
		position: 0,
		link: "",
		errorPatterns: [],
	});
