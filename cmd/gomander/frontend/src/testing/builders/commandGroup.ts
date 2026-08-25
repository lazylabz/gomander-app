import type { Command, CommandGroup } from "@/contracts/types.ts";

export type CommandGroupBuilder = {
	withId: (id: string) => CommandGroupBuilder;
	withProjectId: (projectId: string) => CommandGroupBuilder;
	withName: (name: string) => CommandGroupBuilder;
	withPosition: (position: number) => CommandGroupBuilder;
	withCommands: (...commands: Command[]) => CommandGroupBuilder;
	build: () => CommandGroup;
};

const builder = (data: CommandGroup): CommandGroupBuilder => ({
	withId: (id) => builder({ ...data, id }),
	withProjectId: (projectId) => builder({ ...data, projectId }),
	withName: (name) => builder({ ...data, name }),
	withPosition: (position) => builder({ ...data, position }),
	withCommands: (...commands) => builder({ ...data, commands }),
	build: () => ({ ...data, commands: [...data.commands] }),
});

export const newCommandGroupBuilder = (): CommandGroupBuilder =>
	builder({
		id: crypto.randomUUID(),
		projectId: crypto.randomUUID(),
		name: "Test Command Group",
		position: 0,
		commands: [],
	});
