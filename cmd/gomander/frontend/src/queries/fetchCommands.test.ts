import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { commandStore } from "@/store/commandStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { newCommandBuilder } from "@/testing/builders/command.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

describe("fetchCommands", () => {
	const sut = fetchCommands;

	beforeEach(() => {
		commandStore.setState({ commands: [], commandsStatus: {} });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should load the commands held by the backend into the store", async () => {
		// Arrange
		const firstCommand = newCommandBuilder().withId("cmd-1").build();
		const secondCommand = newCommandBuilder().withId("cmd-2").build();
		installInMemoryBackend({ commands: [firstCommand, secondCommand] });

		// Act
		await sut();

		// Assert
		expect(commandStore.getState().commands).toEqual([
			firstCommand,
			secondCommand,
		]);
	});

	it("Should mark every fetched command as idle", async () => {
		// Arrange
		installInMemoryBackend({
			commands: [newCommandBuilder().withId("cmd-1").build()],
		});

		// Act
		await sut();

		// Assert
		expect(commandStore.getState().commandsStatus).toEqual({
			"cmd-1": CommandStatus.IDLE,
		});
	});

	it("Should keep the status of a command that is already known", async () => {
		// Arrange
		installInMemoryBackend({
			commands: [newCommandBuilder().withId("cmd-1").build()],
		});
		commandStore.setState({
			commandsStatus: { "cmd-1": CommandStatus.RUNNING },
		});

		// Act
		await sut();

		// Assert
		expect(commandStore.getState().commandsStatus).toEqual({
			"cmd-1": CommandStatus.RUNNING,
		});
	});

	it("Should forget the status of a command the backend no longer has", async () => {
		// Arrange
		installInMemoryBackend({
			commands: [newCommandBuilder().withId("cmd-1").build()],
		});
		commandStore.setState({
			commandsStatus: {
				"cmd-1": CommandStatus.IDLE,
				"deleted-cmd": CommandStatus.RUNNING,
			},
		});

		// Act
		await sut();

		// Assert
		expect(commandStore.getState().commandsStatus).toEqual({
			"cmd-1": CommandStatus.IDLE,
		});
	});
});
