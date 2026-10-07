import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { commandActions } from "@/components/ActionPalette/actions/commandActions.ts";
import { resetBackendServices } from "@/contracts/service.ts";
import { commandStore } from "@/store/commandStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { resetStores } from "@/testing/stores.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

describe("commandActions", () => {
	const sut = commandActions;

	beforeEach(async () => {
		await installTranslations();
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should offer to run an idle command and to stop a running one", () => {
		// Arrange
		const idle = new CommandBuilder().withId("idle").build();
		const running = new CommandBuilder().withId("running").build();

		// Act
		const actions = sut([idle, running], {
			idle: CommandStatus.IDLE,
			running: CommandStatus.RUNNING,
		});

		// Assert
		expect(actions.map((action) => action.id)).toEqual([
			"command.start.idle",
			"command.stop.running",
		]);
	});

	it("Should offer nothing for a command whose status is not known yet", () => {
		// Arrange
		const command = new CommandBuilder().build();

		// Act
		const actions = sut([command], {});

		// Assert
		expect(actions).toEqual([]);
	});

	it("Should let the user find a command by its name or its command line", () => {
		// Arrange
		const command = new CommandBuilder()
			.withName("api")
			.withCommand("go run ./cmd/api")
			.build();

		// Act
		const [action] = sut([command], { [command.id]: CommandStatus.IDLE });

		// Assert
		expect(action.keywords).toEqual(["api", "go run ./cmd/api"]);
	});

	it("Should run the command and show its output", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		const command = new CommandBuilder().withId("cmd-1").build();
		const [action] = sut([command], { "cmd-1": CommandStatus.IDLE });

		// Act
		await action.perform();

		// Assert
		expect(backend.state.runningCommandIds).toEqual(["cmd-1"]);
		expect(commandStore.getState().activeCommandId).toBe("cmd-1");
	});

	it("Should stop the command and show its output", async () => {
		// Arrange
		const backend = installInMemoryBackend({ runningCommandIds: ["cmd-1"] });
		const command = new CommandBuilder().withId("cmd-1").build();
		const [action] = sut([command], { "cmd-1": CommandStatus.RUNNING });

		// Act
		await action.perform();

		// Assert
		expect(backend.state.runningCommandIds).toEqual([]);
		expect(commandStore.getState().activeCommandId).toBe("cmd-1");
	});
});
