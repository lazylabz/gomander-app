import type { Project } from "@/contracts/types.ts";

export type ProjectBuilder = {
	withId: (id: string) => ProjectBuilder;
	withName: (name: string) => ProjectBuilder;
	withWorkingDirectory: (workingDirectory: string) => ProjectBuilder;
	build: () => Project;
};

const builder = (data: Project): ProjectBuilder => ({
	withId: (id) => builder({ ...data, id }),
	withName: (name) => builder({ ...data, name }),
	withWorkingDirectory: (workingDirectory) =>
		builder({ ...data, workingDirectory }),
	build: () => ({ ...data }),
});

export const newProjectBuilder = (): ProjectBuilder =>
	builder({
		id: crypto.randomUUID(),
		name: "Test Project",
		workingDirectory: "/app",
	});
