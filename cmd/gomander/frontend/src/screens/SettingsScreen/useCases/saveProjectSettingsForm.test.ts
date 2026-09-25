import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { saveProjectSettingsForm } from "@/screens/SettingsScreen/useCases/saveProjectSettingsForm.ts";
import { projectStore } from "@/store/projectStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { ProjectBuilder } from "@/testing/builders/project.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { resetStores } from "@/testing/stores.ts";

describe("saveProjectSettingsForm", () => {
	const sut = saveProjectSettingsForm;

	const toastSuccess = vi.spyOn(toast, "success");

	const project = new ProjectBuilder()
		.withId("project-1")
		.withName("Old name")
		.withWorkingDirectory("/old")
		.build();

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should save the new name and working directory of the opened project", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			projects: [project],
			currentProject: project,
		});
		projectStore.setState({ projectInfo: project });
		const edited = {
			id: "project-1",
			name: "New name",
			workingDirectory: "/new",
		};

		// Act
		await sut({ name: "New name", baseWorkingDirectory: "/new" });

		// Assert
		expect(backend.state.projects).toEqual([edited]);
		expect(projectStore.getState().projectInfo).toEqual(edited);
	});

	it("Should save nothing when no project is open", async () => {
		// Arrange
		const backend = installInMemoryBackend({ projects: [project] });

		// Act
		await sut({ name: "New name", baseWorkingDirectory: "/new" });

		// Assert
		expect(backend.state.projects).toEqual([project]);
		expect(projectStore.getState().projectInfo).toBeNull();
		expect(toastSuccess).not.toHaveBeenCalled();
	});
});
