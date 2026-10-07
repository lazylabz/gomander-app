import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { commandGroupActions } from "@/components/ActionPalette/actions/commandGroupActions.ts";
import { resetBackendServices } from "@/contracts/service.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

describe("commandGroupActions", () => {
	const sut = commandGroupActions;

	const api = new CommandBuilder().withId("api").build();
	const web = new CommandBuilder().withId("web").build();
	const group = new CommandGroupBuilder()
		.withId("group")
		.withName("stack")
		.withCommands(api, web)
		.build();

	beforeEach(async () => {
		await installTranslations();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should offer only to run a group with nothing running", () => {
		// Act
		const actions = sut([group], [api, web], {
			api: CommandStatus.IDLE,
			web: CommandStatus.IDLE,
		});

		// Assert
		expect(actions.map((action) => action.id)).toEqual([
			"commandGroup.run.group",
		]);
	});

	it("Should offer only to stop a group with everything running", () => {
		// Act
		const actions = sut([group], [api, web], {
			api: CommandStatus.RUNNING,
			web: CommandStatus.RUNNING,
		});

		// Assert
		expect(actions.map((action) => action.id)).toEqual([
			"commandGroup.stop.group",
		]);
	});

	it("Should offer to run and to stop a partly running group", () => {
		// Act
		const actions = sut([group], [api, web], {
			api: CommandStatus.RUNNING,
			web: CommandStatus.IDLE,
		});

		// Assert
		expect(actions.map((action) => action.id)).toEqual([
			"commandGroup.run.group",
			"commandGroup.stop.group",
		]);
	});

	it("Should offer nothing for a group whose commands are not loaded", () => {
		// Act
		const actions = sut([group], [], {});

		// Assert
		expect(actions).toEqual([]);
	});

	it("Should run the group", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commands: [api, web],
			commandGroups: [group],
		});
		const [action] = sut([group], [api, web], {
			api: CommandStatus.IDLE,
			web: CommandStatus.IDLE,
		});

		// Act
		await action.perform();

		// Assert
		expect(backend.state.runningGroupIds).toEqual(["group"]);
	});

	it("Should stop the group", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commands: [api, web],
			commandGroups: [group],
			runningGroupIds: ["group"],
		});
		const [action] = sut([group], [api, web], {
			api: CommandStatus.RUNNING,
			web: CommandStatus.RUNNING,
		});

		// Act
		await action.perform();

		// Assert
		expect(backend.state.runningGroupIds).toEqual([]);
	});
});
