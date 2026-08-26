import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { projectStore } from "@/store/projectStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { ProjectBuilder } from "@/testing/builders/project.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { editOpenedProject } from "@/useCases/project/editOpenedProject.ts";

describe("editOpenedProject", () => {
	const sut = editOpenedProject;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	const project = new ProjectBuilder()
		.withId("project-1")
		.withName("Old name")
		.build();
	const renamed = { ...project, name: "New name" };

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		projectStore.setState({
			projectInfo: project,
			availableProjects: [project],
		});
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should save the project and refresh the opened project and the available ones", async () => {
		// Arrange
		installInMemoryBackend({ projects: [project], currentProject: project });

		// Act
		const succeeded = await sut(renamed);

		// Assert
		expect(succeeded).toBe(true);
		expect(projectStore.getState().projectInfo).toEqual(renamed);
		expect(projectStore.getState().availableProjects).toEqual([renamed]);
		expect(toastSuccess).toHaveBeenCalledWith(
			"toast.settings.projectSaveSuccess",
		);
	});

	it("Should notify the user when the backend rejects the edit", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			projects: [project],
			currentProject: project,
		});
		backend.data.editProject = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut(renamed);

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith(
			"toast.settings.projectSaveFailed: boom",
		);
		expect(projectStore.getState().projectInfo).toEqual(project);
	});
});
