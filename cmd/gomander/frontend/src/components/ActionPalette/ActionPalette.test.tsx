import { screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { ActionPalette } from "@/components/ActionPalette/ActionPalette.tsx";
import { resetBackendServices } from "@/contracts/service.ts";
import type { Command } from "@/contracts/types.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import { resetStores } from "@/testing/stores.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

describe("ActionPalette", () => {
	const seed = (commands: Command[], status: CommandStatus) => {
		const backend = installInMemoryBackend({ commands });
		commandStore.setState({
			commands,
			commandsStatus: Object.fromEntries(commands.map((c) => [c.id, status])),
		});
		return backend;
	};

	beforeEach(async () => {
		await installTranslations();
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should open on Mod+K and close on it again", async () => {
		// Arrange
		seed([new CommandBuilder().build()], CommandStatus.IDLE);
		const { user } = renderWithProviders(<ActionPalette />);

		// Act
		await user.keyboard("{Meta>}k{/Meta}");

		// Assert
		expect(screen.getByRole("dialog")).toBeInTheDocument();
		expect(
			screen.getByRole("option", {
				name: "actionPalette.run actionPalette.kinds.command",
			}),
		).toBeInTheDocument();

		// Act
		await user.keyboard("{Meta>}k{/Meta}");

		// Assert
		expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
	});

	it("Should run the picked action and close", async () => {
		// Arrange
		const command = new CommandBuilder().withId("cmd-1").build();
		const backend = seed([command], CommandStatus.IDLE);
		const { user } = renderWithProviders(<ActionPalette />);
		await user.keyboard("{Meta>}k{/Meta}");

		// Act
		await user.click(
			screen.getByRole("option", {
				name: "actionPalette.run actionPalette.kinds.command",
			}),
		);

		// Assert
		expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
		expect(backend.state.runningCommandIds).toEqual(["cmd-1"]);
	});

	it("Should filter actions by their keywords and not by their id", async () => {
		// Arrange
		seed(
			[
				new CommandBuilder().withId("abc").withName("api").build(),
				new CommandBuilder().withId("xyz").withName("web").build(),
			],
			CommandStatus.IDLE,
		);
		const { user } = renderWithProviders(<ActionPalette />);
		await user.keyboard("{Meta>}k{/Meta}");

		// Act
		await user.keyboard("api");

		// Assert
		expect(screen.getAllByRole("option")).toHaveLength(1);

		// Act
		await user.clear(screen.getByRole("combobox"));
		await user.keyboard("xyz");

		// Assert
		expect(screen.queryAllByRole("option")).toHaveLength(0);
		expect(screen.getByText("actionPalette.empty")).toBeInTheDocument();
	});

	it("Should tell a command group apart from a command of the same name", async () => {
		// Arrange
		const command = new CommandBuilder().withName("stack").build();
		const commandGroup = new CommandGroupBuilder()
			.withName("stack")
			.withCommands(command)
			.build();
		seed([command], CommandStatus.IDLE);
		commandGroupStore.setState({ commandGroups: [commandGroup] });
		const { user } = renderWithProviders(<ActionPalette />);

		// Act
		await user.keyboard("{Meta>}k{/Meta}");

		// Assert
		expect(
			screen.getByRole("option", {
				name: "actionPalette.run actionPalette.kinds.command",
			}),
		).toBeInTheDocument();
		expect(
			screen.getByRole("option", {
				name: "actionPalette.run actionPalette.kinds.commandGroup",
			}),
		).toBeInTheDocument();
	});
});
