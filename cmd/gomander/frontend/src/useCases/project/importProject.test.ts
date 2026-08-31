import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { projectStore } from "@/store/projectStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { ProjectBlueprintBuilder } from "@/testing/builders/projectBlueprint.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { importProject } from "@/useCases/project/importProject.ts";

describe("importProject", () => {
	const sut = importProject;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	const projectToImport = new ProjectBlueprintBuilder()
		.withName("Exported project")
		.build();

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		projectStore.setState({ availableProjects: [] });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should import the project and refresh the available projects", async () => {
		// Arrange
		const backend = installInMemoryBackend();

		// Act
		const succeeded = await sut(projectToImport, "A new name", "/new/dir");

		// Assert
		expect(succeeded).toBe(true);
		expect(backend.state.importedProjects).toEqual([projectToImport]);
		expect(projectStore.getState().availableProjects).toEqual([
			expect.objectContaining({
				name: "A new name",
				workingDirectory: "/new/dir",
			}),
		]);
		expect(toastSuccess).toHaveBeenCalledWith("toast.project.importSuccess");
	});

	it("Should notify the user when the backend rejects the import", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.importProject = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut(projectToImport, "A new name", "/new/dir");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith("toast.project.importFailed: boom");
		expect(projectStore.getState().availableProjects).toEqual([]);
	});
});
