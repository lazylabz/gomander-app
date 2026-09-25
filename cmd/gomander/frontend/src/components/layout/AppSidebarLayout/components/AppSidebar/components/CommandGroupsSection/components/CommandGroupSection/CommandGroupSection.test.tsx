import { DndContext } from "@dnd-kit/core";
import { SortableContext } from "@dnd-kit/sortable";
import { screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { CommandGroupSection } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/CommandGroupsSection/components/CommandGroupSection/CommandGroupSection.tsx";
import type { InMemoryBackend } from "@/contracts/adapters/inMemory.ts";
import { resetBackendServices } from "@/contracts/service.ts";
import type { Command, CommandGroup } from "@/contracts/types.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
import {
	forgetSidebarSection,
	setSidebarSectionOpen,
} from "@/store/sidebarSections.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import { resetStores } from "@/testing/stores.ts";
import {
	installRecordingTerminals,
	resetTerminals,
} from "@/testing/terminals.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

describe("CommandGroupSection", () => {
	const startEditingCommandGroup = vi.fn();

	let backend: InMemoryBackend;

	const api = new CommandBuilder().withId("cmd-1").withName("api").build();
	const web = new CommandBuilder().withId("cmd-2").withName("web").build();
	const group = new CommandGroupBuilder()
		.withId("group-1")
		.withName("Backend")
		.withCommands(api, web)
		.build();

	const seed = (statuses: Record<string, CommandStatus>) => {
		const commands: Command[] = [api, web];
		const commandGroups: CommandGroup[] = [group];
		backend = installInMemoryBackend({ commands, commandGroups });
		commandStore.setState({ commands, commandsStatus: statuses });
		commandGroupStore.setState({ commandGroups });
	};

	const render = ({ isReorderingGroups = false } = {}) =>
		renderWithProviders(
			<DndContext>
				<SortableContext items={[group.id]}>
					<CommandGroupSection
						commandGroup={group}
						startEditingCommandGroup={startEditingCommandGroup}
						isReorderingGroups={isReorderingGroups}
					/>
				</SortableContext>
			</DndContext>,
		);

	const allIdle = {
		"cmd-1": CommandStatus.IDLE,
		"cmd-2": CommandStatus.IDLE,
	};
	const oneRunning = {
		"cmd-1": CommandStatus.RUNNING,
		"cmd-2": CommandStatus.IDLE,
	};
	const allRunning = {
		"cmd-1": CommandStatus.RUNNING,
		"cmd-2": CommandStatus.RUNNING,
	};

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		installRecordingTerminals();
		resetStores();
	});

	afterEach(() => {
		// The sections store is a module singleton persisted to storage, so an
		// opened section would leak into the next test.
		forgetSidebarSection(group.id);
		resetBackendServices();
		resetTerminals();
	});

	it("Should not show a running count when no command runs", () => {
		// Arrange
		seed(allIdle);

		// Act
		render();

		// Assert
		expect(screen.queryByText("(0/2)")).not.toBeInTheDocument();
	});

	it("Should show how many of its commands run when some do", () => {
		// Arrange
		seed(oneRunning);

		// Act
		render();

		// Assert
		expect(screen.getByText("(1/2)")).toBeInTheDocument();
	});

	it("Should offer to run but not to stop the group when every command is idle", () => {
		// Arrange
		seed(allIdle);

		// Act
		render();

		// Assert
		expect(
			screen.getByRole("button", { name: "sidebar.commandGroups.run" }),
		).toBeInTheDocument();
		expect(
			screen.queryByRole("button", { name: "sidebar.commandGroups.stop" }),
		).not.toBeInTheDocument();
	});

	it("Should offer to stop but not to run the group when every command runs", () => {
		// Arrange
		seed(allRunning);

		// Act
		render();

		// Assert
		expect(
			screen.getByRole("button", { name: "sidebar.commandGroups.stop" }),
		).toBeInTheDocument();
		expect(
			screen.queryByRole("button", { name: "sidebar.commandGroups.run" }),
		).not.toBeInTheDocument();
	});

	it("Should run the group without expanding it when run is clicked", async () => {
		// Arrange
		seed(allIdle);
		const { user } = render();

		// Act
		await user.click(
			screen.getByRole("button", { name: "sidebar.commandGroups.run" }),
		);

		// Assert
		expect(backend.state.runningGroupIds).toEqual(["group-1"]);
		expect(screen.queryByText("api")).not.toBeInTheDocument();
	});

	it("Should stop the group when stop is clicked", async () => {
		// Arrange
		seed(allRunning);
		backend.state.runningGroupIds = ["group-1"];
		const { user } = render();

		// Act
		await user.click(
			screen.getByRole("button", { name: "sidebar.commandGroups.stop" }),
		);

		// Assert
		expect(backend.state.runningGroupIds).toEqual([]);
	});

	it("Should show the group's commands once its header is clicked", async () => {
		// Arrange
		seed(allIdle);
		const { user } = render();

		// Act
		await user.click(screen.getByText("Backend"));

		// Assert
		expect(screen.getByText("api")).toBeInTheDocument();
		expect(screen.getByText("web")).toBeInTheDocument();
	});

	it("Should hide the group's commands when the header of an open group is clicked", async () => {
		// Arrange
		seed(allIdle);
		setSidebarSectionOpen(group.id, true);
		const { user } = render();

		// Act
		await user.click(screen.getByText("Backend"));

		// Assert
		expect(screen.queryByText("api")).not.toBeInTheDocument();
	});

	it("Should hide the commands and the controls of an open group while groups are reordered", () => {
		// Arrange
		seed(oneRunning);
		setSidebarSectionOpen(group.id, true);

		// Act
		render({ isReorderingGroups: true });

		// Assert
		expect(screen.queryByText("api")).not.toBeInTheDocument();
		expect(screen.queryByText("(1/2)")).not.toBeInTheDocument();
		expect(
			screen.queryByRole("button", { name: "sidebar.commandGroups.run" }),
		).not.toBeInTheDocument();
		expect(
			screen.queryByRole("button", { name: "sidebar.commandGroups.stop" }),
		).not.toBeInTheDocument();
	});

	it("Should hand the group to the editor when edit is chosen", async () => {
		// Arrange
		seed(allIdle);
		const { user } = render();
		await user.pointer({
			keys: "[MouseRight]",
			target: screen.getByText("Backend"),
		});

		// Act
		await user.click(screen.getByRole("menuitem", { name: "common.edit" }));

		// Assert
		expect(startEditingCommandGroup).toHaveBeenCalledWith(group);
	});

	it("Should delete the group when delete is chosen", async () => {
		// Arrange
		seed(allIdle);
		const { user } = render();
		await user.pointer({
			keys: "[MouseRight]",
			target: screen.getByText("Backend"),
		});

		// Act
		await user.click(screen.getByRole("menuitem", { name: "common.delete" }));

		// Assert
		await waitFor(() => expect(backend.state.commandGroups).toEqual([]));
	});
});
