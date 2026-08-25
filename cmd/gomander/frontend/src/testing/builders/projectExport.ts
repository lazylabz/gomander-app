import type {
	CommandExport,
	CommandGroupExport,
	ProjectExport,
} from "@/contracts/types.ts";

export type ProjectExportBuilder = {
	withVersion: (version: number) => ProjectExportBuilder;
	withName: (name: string) => ProjectExportBuilder;
	withWorkingDirectory: (workingDirectory: string) => ProjectExportBuilder;
	withCommands: (...commands: CommandExport[]) => ProjectExportBuilder;
	withCommandGroups: (
		...commandGroups: CommandGroupExport[]
	) => ProjectExportBuilder;
	build: () => ProjectExport;
};

const builder = (data: ProjectExport): ProjectExportBuilder => ({
	withVersion: (version) => builder({ ...data, version }),
	withName: (name) => builder({ ...data, name }),
	withWorkingDirectory: (workingDirectory) =>
		builder({ ...data, workingDirectory }),
	withCommands: (...commands) => builder({ ...data, commands }),
	withCommandGroups: (...commandGroups) => builder({ ...data, commandGroups }),
	build: () => ({
		...data,
		commands: [...data.commands],
		commandGroups: [...data.commandGroups],
	}),
});

export const newProjectExportBuilder = (): ProjectExportBuilder =>
	builder({
		version: 1,
		name: "Test Project",
		workingDirectory: "/app",
		commands: [],
		commandGroups: [],
	});
