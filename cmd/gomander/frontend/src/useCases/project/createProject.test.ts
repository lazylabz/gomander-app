import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { projectStore } from "@/store/projectStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { resetStores } from "@/testing/stores.ts";
import { createProject } from "@/useCases/project/createProject.ts";

describe("createProject", () => {
	const sut = createProject;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should create the project and refresh the available projects", async () => {
		// Arrange
		const backend = installInMemoryBackend();

		// Act
		const succeeded = await sut("My project", "/my/dir");

		// Assert
		const expected = expect.objectContaining({
			name: "My project",
			workingDirectory: "/my/dir",
		});
		expect(succeeded).toBe(true);
		expect(backend.state.projects).toEqual([expected]);
		expect(projectStore.getState().availableProjects).toEqual([expected]);
	});

	it("Should notify the user that the project was created", async () => {
		// Arrange
		installInMemoryBackend();

		// Act
		await sut("My project", "/my/dir");

		// Assert
		expect(toastSuccess).toHaveBeenCalledWith("toast.project.createSuccess");
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should still report success when the refresh that follows fails", async () => {
		// Arrange
		const consoleError = vi
			.spyOn(console, "error")
			.mockImplementation(() => {});
		const backend = installInMemoryBackend();
		backend.data.getAvailableProjects = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut("My project", "/my/dir");

		// Assert
		expect(succeeded).toBe(true);
		expect(toastSuccess).toHaveBeenCalledWith("toast.project.createSuccess");
		expect(toastError).not.toHaveBeenCalled();

		consoleError.mockRestore();
	});

	it("Should create no project and notify the user when the backend rejects the creation", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.createProject = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut("My project", "/my/dir");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith("toast.project.createFailed: boom");
		expect(toastSuccess).not.toHaveBeenCalled();
		expect(backend.state.projects).toEqual([]);
		expect(projectStore.getState().availableProjects).toEqual([]);
	});
});
