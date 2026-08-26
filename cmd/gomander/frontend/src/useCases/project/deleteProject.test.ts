import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { projectStore } from "@/store/projectStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { ProjectBuilder } from "@/testing/builders/project.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { deleteProject } from "@/useCases/project/deleteProject.ts";

describe("deleteProject", () => {
	const sut = deleteProject;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	const firstProject = new ProjectBuilder().withId("project-1").build();
	const secondProject = new ProjectBuilder().withId("project-2").build();

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		projectStore.setState({
			availableProjects: [firstProject, secondProject],
		});
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should delete the project and refresh the available projects", async () => {
		// Arrange
		installInMemoryBackend({ projects: [firstProject, secondProject] });

		// Act
		const succeeded = await sut("project-1");

		// Assert
		expect(succeeded).toBe(true);
		expect(projectStore.getState().availableProjects).toEqual([secondProject]);
		expect(toastSuccess).toHaveBeenCalledWith("toast.project.deleteSuccess");
	});

	it("Should keep the project and notify the user when the backend rejects the deletion", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			projects: [firstProject, secondProject],
		});
		backend.data.deleteProject = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut("project-1");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith("toast.project.deleteFailed: boom");
		expect(projectStore.getState().availableProjects).toEqual([
			firstProject,
			secondProject,
		]);
	});
});
