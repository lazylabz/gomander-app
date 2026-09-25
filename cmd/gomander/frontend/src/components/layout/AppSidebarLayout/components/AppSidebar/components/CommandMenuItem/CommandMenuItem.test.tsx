import { DndContext } from "@dnd-kit/core";
import { SortableContext } from "@dnd-kit/sortable";
import { screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { CommandMenuItem } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/CommandMenuItem/CommandMenuItem.tsx";
import { sidebarContext } from "@/components/layout/AppSidebarLayout/components/AppSidebar/contexts/sidebarContext.tsx";
import type { InMemoryBackend } from "@/contracts/adapters/inMemory.ts";
import { resetBackendServices } from "@/contracts/service.ts";
import type { Command, CommandGroup } from "@/contracts/types.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
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

describe("CommandMenuItem", () => {
	const startEditingCommand = vi.fn();

	let backend: InMemoryBackend;

	const render = (command: Command, insideGroupId?: string) =>
		renderWithProviders(
			<sidebarContext.Provider value={{ startEditingCommand }}>
				<DndContext>
					<SortableContext items={[command.id]}>
						<CommandMenuItem command={command} insideGroupId={insideGroupId} />
					</SortableContext>
				</DndContext>
			</sidebarContext.Provider>,
		);

	const seed = (
		command: Command,
		status: CommandStatus,
		commandGroups: CommandGroup[] = [],
	) => {
		backend = installInMemoryBackend({ commands: [command], commandGroups });
		commandStore.setState({
			commands: [command],
			commandsStatus: { [command.id]: status },
		});
		commandGroupStore.setState({ commandGroups });
	};

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		installRecordingTerminals();
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
		resetTerminals();
	});

	it("Should offer to run an idle command but not to stop it", () => {
		// Arrange
		const command = new CommandBuilder().build();
		seed(command, CommandStatus.IDLE);

		// Act
		render(command);

		// Assert
		expect(
			screen.getByRole("button", { name: "sidebar.commands.run" }),
		).toBeInTheDocument();
		expect(
			screen.queryByRole("button", { name: "sidebar.commands.stop" }),
		).not.toBeInTheDocument();
	});

	it("Should offer to stop a running command but not to run it", () => {
		// Arrange
		const command = new CommandBuilder().build();
		seed(command, CommandStatus.RUNNING);

		// Act
		render(command);

		// Assert
		expect(
			screen.getByRole("button", { name: "sidebar.commands.stop" }),
		).toBeInTheDocument();
		expect(
			screen.queryByRole("button", { name: "sidebar.commands.run" }),
		).not.toBeInTheDocument();
	});

	it("Should start the command and make it active when run is clicked", async () => {
		// Arrange
		const command = new CommandBuilder().withId("cmd-1").build();
		seed(command, CommandStatus.IDLE);
		const { user } = render(command);

		// Act
		await user.click(
			screen.getByRole("button", { name: "sidebar.commands.run" }),
		);

		// Assert
		expect(backend.state.runningCommandIds).toEqual(["cmd-1"]);
		expect(commandStore.getState().activeCommandId).toBe("cmd-1");
	});

	it("Should stop the command when stop is clicked", async () => {
		// Arrange
		const command = new CommandBuilder().withId("cmd-1").build();
		seed(command, CommandStatus.RUNNING);
		backend.state.runningCommandIds = ["cmd-1"];
		const { user } = render(command);

		// Act
		await user.click(
			screen.getByRole("button", { name: "sidebar.commands.stop" }),
		);

		// Assert
		expect(backend.state.runningCommandIds).toEqual([]);
	});

	it("Should make the command active and clear its error when the row is clicked", async () => {
		// Arrange
		const command = new CommandBuilder()
			.withId("cmd-1")
			.withName("api")
			.build();
		seed(command, CommandStatus.IDLE);
		commandStore.setState({ commandIdsWithErrors: ["cmd-1"] });
		const { user } = render(command);

		// Act
		await user.click(screen.getByText("api"));

		// Assert
		expect(commandStore.getState().activeCommandId).toBe("cmd-1");
		expect(commandStore.getState().commandIdsWithErrors).toEqual([]);
	});

	it("Should not offer a link control when the command has no link", () => {
		// Arrange
		const command = new CommandBuilder().withLink("").build();
		seed(command, CommandStatus.IDLE);

		// Act
		render(command);

		// Assert
		expect(
			screen.queryByRole("button", { name: "sidebar.commands.openLink" }),
		).not.toBeInTheDocument();
	});

	it("Should open the command's link when the link control is clicked", async () => {
		// Arrange
		const command = new CommandBuilder()
			.withLink("http://localhost:3000")
			.build();
		seed(command, CommandStatus.IDLE);
		const { user } = render(command);

		// Act
		await user.click(
			screen.getByRole("button", { name: "sidebar.commands.openLink" }),
		);

		// Assert
		expect(backend.state.openedUrls).toEqual(["http://localhost:3000"]);
	});

	it("Should disable editing, duplicating and deleting while the command runs", async () => {
		// Arrange
		const command = new CommandBuilder().withName("api").build();
		seed(command, CommandStatus.RUNNING);
		const { user } = render(command);

		// Act
		await user.pointer({
			keys: "[MouseRight]",
			target: screen.getByText("api"),
		});

		// Assert
		for (const name of ["common.edit", "common.duplicate", "common.delete"]) {
			expect(screen.getByRole("menuitem", { name })).toHaveAttribute(
				"aria-disabled",
				"true",
			);
		}
	});

	it.each([
		{
			situation: "inside a group",
			insideGroupId: "group-1",
			offered: "sidebar.commands.removeFromGroup",
			withheld: "common.delete",
		},
		{
			situation: "outside a group",
			insideGroupId: undefined,
			offered: "common.delete",
			withheld: "sidebar.commands.removeFromGroup",
		},
	])("Should offer $offered but not $withheld $situation", async ({
		insideGroupId,
		offered,
		withheld,
	}) => {
		// Arrange
		const command = new CommandBuilder().withName("api").build();
		const group = new CommandGroupBuilder()
			.withId("group-1")
			.withCommands(command)
			.build();
		seed(command, CommandStatus.IDLE, [group]);
		const { user } = render(command, insideGroupId);

		// Act
		await user.pointer({
			keys: "[MouseRight]",
			target: screen.getByText("api"),
		});

		// Assert
		expect(screen.getByRole("menuitem", { name: offered })).toBeInTheDocument();
		expect(
			screen.queryByRole("menuitem", { name: withheld }),
		).not.toBeInTheDocument();
	});

	it("Should hand the command to the editor when edit is chosen", async () => {
		// Arrange
		const command = new CommandBuilder().withName("api").build();
		seed(command, CommandStatus.IDLE);
		const { user } = render(command);
		await user.pointer({
			keys: "[MouseRight]",
			target: screen.getByText("api"),
		});

		// Act
		await user.click(screen.getByRole("menuitem", { name: "common.edit" }));

		// Assert
		expect(startEditingCommand).toHaveBeenCalledWith(command);
	});

	it("Should duplicate the command into its group when duplicate is chosen inside a group", async () => {
		// Arrange
		const command = new CommandBuilder()
			.withId("cmd-1")
			.withName("api")
			.build();
		const group = new CommandGroupBuilder()
			.withId("group-1")
			.withCommands(command)
			.build();
		seed(command, CommandStatus.IDLE, [group]);
		const { user } = render(command, "group-1");
		await user.pointer({
			keys: "[MouseRight]",
			target: screen.getByText("api"),
		});

		// Act
		await user.click(
			screen.getByRole("menuitem", { name: "common.duplicate" }),
		);

		// Assert
		await waitFor(() =>
			expect(backend.state.commandGroups[0].commandIds).toEqual([
				"cmd-1",
				"cmd-1-copy",
			]),
		);
	});

	it("Should delete the command when delete is chosen outside a group", async () => {
		// Arrange
		const command = new CommandBuilder()
			.withId("cmd-1")
			.withName("api")
			.build();
		seed(command, CommandStatus.IDLE);
		const { user } = render(command);
		await user.pointer({
			keys: "[MouseRight]",
			target: screen.getByText("api"),
		});

		// Act
		await user.click(screen.getByRole("menuitem", { name: "common.delete" }));

		// Assert
		await waitFor(() => expect(backend.state.commands).toEqual([]));
	});

	it("Should remove the command from its group instead of deleting it inside a group", async () => {
		// Arrange
		const command = new CommandBuilder()
			.withId("cmd-1")
			.withName("api")
			.build();
		const group = new CommandGroupBuilder()
			.withId("group-1")
			.withCommands(command, new CommandBuilder().withId("cmd-2").build())
			.build();
		seed(command, CommandStatus.IDLE, [group]);
		const { user } = render(command, "group-1");
		await user.pointer({
			keys: "[MouseRight]",
			target: screen.getByText("api"),
		});

		// Act
		await user.click(
			screen.getByRole("menuitem", {
				name: "sidebar.commands.removeFromGroup",
			}),
		);

		// Assert
		await waitFor(() =>
			expect(backend.state.commandGroups[0].commandIds).toEqual(["cmd-2"]),
		);
		expect(backend.state.commands).toEqual([command]);
	});
});
