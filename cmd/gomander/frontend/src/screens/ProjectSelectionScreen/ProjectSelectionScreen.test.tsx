import { screen, waitFor, within } from "@testing-library/react";
import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import type { InMemoryBackend } from "@/contracts/adapters/inMemory.ts";
import { resetBackendServices } from "@/contracts/service.ts";
import { ScreenRoutes } from "@/routes.ts";
import { ProjectSelectionScreen } from "@/screens/ProjectSelectionScreen/ProjectSelectionScreen.tsx";
import { projectStore } from "@/store/projectStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { ProjectBuilder } from "@/testing/builders/project.ts";
import { ProjectBlueprintBuilder } from "@/testing/builders/projectBlueprint.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import { resetStores } from "@/testing/stores.ts";

describe("ProjectSelectionScreen", () => {
	const renderSut = () => renderWithProviders(<ProjectSelectionScreen />);

	const openProjectOptions = async (
		user: ReturnType<typeof renderSut>["user"],
		projectName: string,
	) => {
		const card = (await screen.findByRole("button", { name: projectName }))
			.parentElement as HTMLElement;
		await user.click(
			within(card).getByRole("button", {
				name: "projectSelection.projectOptions",
			}),
		);
	};

	let backend: InMemoryBackend;

	beforeEach(async () => {
		vi.restoreAllMocks();
		await installTranslations();
		resetStores();
		backend = installInMemoryBackend();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should list the available projects fetched on mount", async () => {
		// Arrange
		backend.state.projects = [
			new ProjectBuilder().withId("p-1").withName("Alpha").build(),
			new ProjectBuilder().withId("p-2").withName("Beta").build(),
		];

		// Act
		renderSut();

		// Assert
		expect(await screen.findByText("Alpha")).toBeInTheDocument();
		expect(screen.getByText("Beta")).toBeInTheDocument();
		expect(
			screen.getByRole("heading", { name: "projectSelection.openTitle" }),
		).toBeInTheDocument();
		expect(
			screen.queryByText("projectSelection.emptyState"),
		).not.toBeInTheDocument();
	});

	it("Should show the empty state when there are no projects", async () => {
		// Arrange / Act
		renderSut();

		// Assert
		expect(
			await screen.findByText("projectSelection.emptyState"),
		).toBeInTheDocument();
		expect(
			screen.getByRole("heading", { name: "projectSelection.welcome" }),
		).toBeInTheDocument();
	});

	it("Should go to the logs when a project is open", async () => {
		// Arrange
		projectStore.getState().setProjectInfo(new ProjectBuilder().build());

		// Act
		const { location } = renderSut();

		// Assert
		await waitFor(() => expect(location().pathname).toBe(ScreenRoutes.Logs));
	});

	it("Should stay on the selection while no project is open", async () => {
		// Arrange / Act
		const { location } = renderSut();

		// Assert
		await screen.findByText("projectSelection.emptyState");
		expect(location().pathname).toBe(ScreenRoutes.ProjectSelection);
	});

	it("Should open the project whose card is clicked and go to the logs", async () => {
		// Arrange
		backend.state.projects = [
			new ProjectBuilder().withId("p-1").withName("Alpha").build(),
		];
		const { user, location } = renderSut();

		// Act
		await user.click(await screen.findByRole("button", { name: "Alpha" }));

		// Assert
		expect(backend.state.currentProject?.id).toBe("p-1");
		await waitFor(() => expect(location().pathname).toBe(ScreenRoutes.Logs));
	});

	it("Should open the import modal with the picked exported project", async () => {
		// Arrange
		backend.state.projectToImport = new ProjectBlueprintBuilder()
			.withName("Exported")
			.build();
		backend.state.packageJsonProjectToImport = new ProjectBlueprintBuilder()
			.withName("From package.json")
			.build();
		const { user } = renderSut();

		// Act
		await user.click(
			screen.getByRole("button", { name: "projectSelection.importButton" }),
		);

		// Assert
		const dialog = await screen.findByRole("dialog", {
			name: "modal.importProject.title",
		});
		expect(within(dialog).getByDisplayValue("Exported")).toBeInTheDocument();
	});

	it("Should open the import modal with the project picked from a package.json", async () => {
		// Arrange
		backend.state.projectToImport = new ProjectBlueprintBuilder()
			.withName("Exported")
			.build();
		backend.state.packageJsonProjectToImport = new ProjectBlueprintBuilder()
			.withName("From package.json")
			.build();
		const { user } = renderSut();

		// Act
		await user.click(
			screen.getByRole("button", { name: "projectSelection.moreOptions" }),
		);
		await user.click(
			await screen.findByRole("menuitem", {
				name: "projectSelection.importPackageJson",
			}),
		);

		// Assert
		const dialog = await screen.findByRole("dialog", {
			name: "modal.importProject.title",
		});
		expect(
			within(dialog).getByDisplayValue("From package.json"),
		).toBeInTheDocument();
	});

	it("Should open no modal when picking the project to import fails", async () => {
		// Arrange
		const toastError = vi.spyOn(toast, "error");
		backend.data.getProjectToImport = async () => {
			throw new Error("boom");
		};
		const { user } = renderSut();

		// Act
		await user.click(
			screen.getByRole("button", { name: "projectSelection.importButton" }),
		);

		// Assert
		await waitFor(() => expect(toastError).toHaveBeenCalled());
		expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
	});

	it("Should ask for confirmation before deleting a project", async () => {
		// Arrange
		backend.state.projects = [
			new ProjectBuilder().withId("p-1").withName("Alpha").build(),
		];
		const { user } = renderSut();

		// Act
		await openProjectOptions(user, "Alpha");
		await user.click(
			await screen.findByRole("menuitem", { name: "common.delete" }),
		);

		// Assert
		expect(
			await screen.findByRole("alertdialog", {
				name: "modal.deleteProject.title",
			}),
		).toBeInTheDocument();
		expect(backend.state.projects).toHaveLength(1);
	});

	it("Should remove the project from the list once the deletion is confirmed", async () => {
		// Arrange
		backend.state.projects = [
			new ProjectBuilder().withId("p-1").withName("Alpha").build(),
			new ProjectBuilder().withId("p-2").withName("Beta").build(),
		];
		const { user } = renderSut();
		await openProjectOptions(user, "Alpha");
		await user.click(
			await screen.findByRole("menuitem", { name: "common.delete" }),
		);
		const confirmation = await screen.findByRole("alertdialog");

		// Act
		await user.click(
			within(confirmation).getByRole("button", { name: "common.delete" }),
		);

		// Assert
		await waitFor(() =>
			expect(screen.queryByText("Alpha")).not.toBeInTheDocument(),
		);
		expect(screen.getByText("Beta")).toBeInTheDocument();
		expect(screen.queryByRole("alertdialog")).not.toBeInTheDocument();
		expect(backend.state.projects.map((p) => p.id)).toEqual(["p-2"]);
	});

	it("Should keep the project when the deletion is cancelled", async () => {
		// Arrange
		backend.state.projects = [
			new ProjectBuilder().withId("p-1").withName("Alpha").build(),
		];
		const { user } = renderSut();
		await openProjectOptions(user, "Alpha");
		await user.click(
			await screen.findByRole("menuitem", { name: "common.delete" }),
		);
		const confirmation = await screen.findByRole("alertdialog");

		// Act
		await user.click(
			within(confirmation).getByRole("button", { name: "common.cancel" }),
		);

		// Assert
		await waitFor(() =>
			expect(screen.queryByRole("alertdialog")).not.toBeInTheDocument(),
		);
		expect(screen.getByText("Alpha")).toBeInTheDocument();
		expect(backend.state.projects).toHaveLength(1);
	});

	it("Should open the create project modal", async () => {
		// Arrange
		const { user } = renderSut();

		// Act
		await user.click(
			screen.getByRole("button", { name: "projectSelection.createButton" }),
		);

		// Assert
		expect(
			await screen.findByRole("dialog", { name: "modal.createProject.title" }),
		).toBeInTheDocument();
	});
});
