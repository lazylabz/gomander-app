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

	describe("when the initial load fails", () => {
		const failProjectsLoad = () => {
			const loadProjects = backend.data.getAvailableProjects;
			backend.data.getAvailableProjects = async () => {
				throw new Error("boom");
			};
			return () => {
				backend.data.getAvailableProjects = loadProjects;
			};
		};

		beforeEach(() => {
			vi.spyOn(console, "error").mockImplementation(() => {});
		});

		it("Should tell the user the load failed instead of showing the empty state", async () => {
			// Arrange
			failProjectsLoad();

			// Act
			renderSut();

			// Assert
			expect(await screen.findByText("common.loadFailed")).toBeInTheDocument();
			expect(
				screen.getByRole("button", { name: "common.retry" }),
			).toBeInTheDocument();
			expect(
				screen.queryByText("projectSelection.emptyState"),
			).not.toBeInTheDocument();
		});

		it("Should list the projects once retrying succeeds", async () => {
			// Arrange
			backend.state.projects = [
				new ProjectBuilder().withId("p-1").withName("Alpha").build(),
			];
			const recoverBackend = failProjectsLoad();
			const { user } = renderSut();
			await screen.findByText("common.loadFailed");
			recoverBackend();

			// Act
			await user.click(screen.getByRole("button", { name: "common.retry" }));

			// Assert
			expect(await screen.findByText("Alpha")).toBeInTheDocument();
			expect(screen.queryByText("common.loadFailed")).not.toBeInTheDocument();
		});

		it("Should keep the error while retrying keeps failing", async () => {
			// Arrange
			failProjectsLoad();
			const { user } = renderSut();
			await screen.findByText("common.loadFailed");

			// Act
			await user.click(screen.getByRole("button", { name: "common.retry" }));

			// Assert
			expect(await screen.findByText("common.loadFailed")).toBeInTheDocument();
			expect(
				screen.queryByText("projectSelection.emptyState"),
			).not.toBeInTheDocument();
		});

		it("Should keep the error, and not the empty state, while a retry is still loading", async () => {
			// Arrange
			failProjectsLoad();
			const { user } = renderSut();
			await screen.findByText("common.loadFailed");
			let finishLoading: () => void = () => {};
			backend.data.getAvailableProjects = () =>
				new Promise((_, reject) => {
					finishLoading = () => reject(new Error("boom"));
				});

			// Act
			await user.click(screen.getByRole("button", { name: "common.retry" }));

			// Assert
			expect(screen.getByText("common.loadFailed")).toBeInTheDocument();
			expect(
				screen.getByRole("button", { name: "common.retry" }),
			).toBeDisabled();
			expect(
				screen.queryByText("projectSelection.emptyState"),
			).not.toBeInTheDocument();
			finishLoading();
			await waitFor(() =>
				expect(
					screen.getByRole("button", { name: "common.retry" }),
				).toBeEnabled(),
			);
		});

		it("Should show neither the error nor the empty state while the first load is pending", async () => {
			// Arrange
			backend.data.getAvailableProjects = () => new Promise(() => {});

			// Act
			renderSut();

			// Assert
			expect(
				screen.queryByText("projectSelection.emptyState"),
			).not.toBeInTheDocument();
			expect(screen.queryByText("common.loadFailed")).not.toBeInTheDocument();
		});

		it("Should tell the user when the project data of an opened project fails to load", async () => {
			// Arrange
			backend.state.projects = [
				new ProjectBuilder().withId("p-1").withName("Alpha").build(),
			];
			const loadCommands = backend.data.getCommands;
			backend.data.getCommands = async () => {
				throw new Error("boom");
			};
			const { user, location } = renderSut();

			// Act
			await user.click(await screen.findByRole("button", { name: "Alpha" }));

			// Assert
			expect(await screen.findByText("common.loadFailed")).toBeInTheDocument();
			expect(location().pathname).toBe(ScreenRoutes.ProjectSelection);
			expect(screen.getByText("Alpha")).toBeInTheDocument();

			backend.data.getCommands = loadCommands;
			await user.click(screen.getByRole("button", { name: "common.retry" }));
			await waitFor(() => expect(location().pathname).toBe(ScreenRoutes.Logs));
		});

		it("Should not let another project be selected while one is still opening", async () => {
			// Arrange
			backend.state.projects = [
				new ProjectBuilder().withId("p-1").withName("Alpha").build(),
				new ProjectBuilder().withId("p-2").withName("Beta").build(),
			];
			const openedProjectIds: string[] = [];
			let failOpening: () => void = () => {};
			backend.data.openProject = (projectId) => {
				openedProjectIds.push(projectId);
				return new Promise((_, reject) => {
					failOpening = () => reject(new Error("boom"));
				});
			};
			const { user } = renderSut();
			await user.click(await screen.findByRole("button", { name: "Alpha" }));

			// Act
			await user.click(screen.getByRole("button", { name: "Beta" }));

			// Assert
			expect(screen.getByRole("button", { name: "Beta" })).toBeDisabled();
			expect(openedProjectIds).toEqual(["p-1"]);
			failOpening();
			await waitFor(() =>
				expect(screen.getByRole("button", { name: "Beta" })).toBeEnabled(),
			);
		});

		it("Should tell the user when opening the project itself fails", async () => {
			// Arrange
			backend.state.projects = [
				new ProjectBuilder().withId("p-1").withName("Alpha").build(),
			];
			backend.data.openProject = async () => {
				throw new Error("boom");
			};
			const { user, location } = renderSut();

			// Act
			await user.click(await screen.findByRole("button", { name: "Alpha" }));

			// Assert
			expect(await screen.findByText("common.loadFailed")).toBeInTheDocument();
			expect(location().pathname).toBe(ScreenRoutes.ProjectSelection);
		});
	});
});
