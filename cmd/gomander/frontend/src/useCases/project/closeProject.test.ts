import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { projectStore } from "@/store/projectStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { ProjectBuilder } from "@/testing/builders/project.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { closeProject } from "@/useCases/project/closeProject.ts";

describe("closeProject", () => {
	const sut = closeProject;

	const toastError = vi.spyOn(toast, "error");

	const project = new ProjectBuilder().withId("project-1").build();

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		projectStore.setState({ projectInfo: project });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should close the project and forget the one that was open", async () => {
		// Arrange
		const backend = installInMemoryBackend({ currentProject: project });

		// Act
		const succeeded = await sut();

		// Assert
		expect(succeeded).toBe(true);
		expect(backend.state.currentProject).toBeNull();
		expect(projectStore.getState().projectInfo).toBeNull();
	});

	it("Should keep the project open and notify the user when the backend rejects the close", async () => {
		// Arrange
		const backend = installInMemoryBackend({ currentProject: project });
		backend.data.closeProject = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut();

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith("toast.project.closeFailed: boom");
		expect(projectStore.getState().projectInfo).toEqual(project);
	});
});
