import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { ProjectBlueprintBuilder } from "@/testing/builders/projectBlueprint.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { pickProjectToImport } from "@/useCases/project/pickProjectToImport.ts";

describe("pickProjectToImport", () => {
	const sut = pickProjectToImport;

	const toastError = vi.spyOn(toast, "error");

	const exportedProject = new ProjectBlueprintBuilder()
		.withName("Exported project")
		.build();
	const packageJsonProject = new ProjectBlueprintBuilder()
		.withName("From package.json")
		.build();

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should return the blueprint of the exported project", async () => {
		// Arrange
		installInMemoryBackend({
			projectToImport: exportedProject,
			packageJsonProjectToImport: packageJsonProject,
		});

		// Act
		const blueprint = await sut("exportedProject");

		// Assert
		expect(blueprint).toEqual(exportedProject);
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should return the blueprint read from the package.json", async () => {
		// Arrange
		installInMemoryBackend({
			projectToImport: exportedProject,
			packageJsonProjectToImport: packageJsonProject,
		});

		// Act
		const blueprint = await sut("packageJson");

		// Assert
		expect(blueprint).toEqual(packageJsonProject);
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should return nothing and notify the user when the pick is rejected", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.getProjectToImport = async () => {
			throw new Error("boom");
		};

		// Act
		const blueprint = await sut("exportedProject");

		// Assert
		expect(blueprint).toBeNull();
		expect(toastError).toHaveBeenCalledWith("toast.project.selectFailed: boom");
	});
});
