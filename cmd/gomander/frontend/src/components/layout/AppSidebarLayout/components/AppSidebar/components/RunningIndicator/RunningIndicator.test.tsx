import { screen } from "@testing-library/react";
import { beforeEach, describe, expect, it } from "vitest";

import { RunningIndicator } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/RunningIndicator/RunningIndicator.tsx";
import { commandStore } from "@/store/commandStore.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { installTranslations, withTranslation } from "@/testing/i18n.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import { resetStores } from "@/testing/stores.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

describe("RunningIndicator", () => {
	beforeEach(async () => {
		await installTranslations();
		resetStores();
	});

	it("Should tell how many loaded commands run", async () => {
		// Arrange
		withTranslation(
			"sidebar.runningIndicator.running_other",
			"{{count}} running",
		);
		commandStore.setState({
			commands: [
				new CommandBuilder().withId("cmd-1").build(),
				new CommandBuilder().withId("cmd-2").build(),
			],
			commandsStatus: {
				"cmd-1": CommandStatus.RUNNING,
				"cmd-2": CommandStatus.RUNNING,
			},
		});
		const { user } = renderWithProviders(<RunningIndicator />);

		// Act
		await user.hover(screen.getByRole("button"));

		// Assert
		expect(await screen.findByRole("tooltip")).toHaveTextContent("2 running");
	});

	// A status can outlive its command (another project's, or one just deleted),
	// so only the loaded commands count.
	it("Should tell that nothing runs when the only running status belongs to a command not loaded", async () => {
		// Arrange
		commandStore.setState({
			commands: [new CommandBuilder().withId("cmd-1").build()],
			commandsStatus: {
				"cmd-1": CommandStatus.IDLE,
				"gone-cmd": CommandStatus.RUNNING,
			},
		});
		const { user } = renderWithProviders(<RunningIndicator />);

		// Act
		await user.hover(screen.getByRole("button"));

		// Assert
		expect(await screen.findByRole("tooltip")).toHaveTextContent(
			"sidebar.runningIndicator.idle",
		);
	});
});
