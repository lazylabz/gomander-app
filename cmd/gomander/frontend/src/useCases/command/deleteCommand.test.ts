import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import type { RecordingTerminals } from "@/commandOutput/adapters/recording.ts";

import {
	appendCommandOutput,
	attachCommandOutput,
	commandOutputTail,
} from "@/commandOutput/commandOutput.ts";
import { resetBackendServices } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { installTranslations } from "@/testing/i18n.ts";
import {
	installRecordingTerminals,
	resetTerminals,
} from "@/testing/terminals.ts";
import { deleteCommand } from "@/useCases/command/deleteCommand.ts";

describe("deleteCommand", () => {
	const sut = deleteCommand;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	let recording: RecordingTerminals;

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		recording = installRecordingTerminals();
		commandStore.setState({
			commands: [],
			commandsStatus: {},
			activeCommandId: null,
		});
		commandGroupStore.setState({ commandGroups: [] });
	});

	afterEach(() => {
		resetBackendServices();
		resetTerminals();
	});

	it("Should delete the command and refresh the commands and the groups", async () => {
		// Arrange
		const command = new CommandBuilder().withId("cmd-1").build();
		const group = new CommandGroupBuilder()
			.withId("group-1")
			.withCommands(command, new CommandBuilder().withId("cmd-2").build())
			.build();
		installInMemoryBackend({ commands: [command], commandGroups: [group] });

		// Act
		const succeeded = await sut("cmd-1");

		// Assert
		expect(succeeded).toBe(true);
		expect(commandStore.getState().commands).toEqual([]);
		expect(commandGroupStore.getState().commandGroups[0].commandIds).toEqual([
			"cmd-2",
		]);
	});

	it("Should clear the active command", async () => {
		// Arrange
		installInMemoryBackend({
			commands: [new CommandBuilder().withId("cmd-1").build()],
		});
		commandStore.setState({ activeCommandId: "cmd-1" });

		// Act
		await sut("cmd-1");

		// Assert
		expect(commandStore.getState().activeCommandId).toBeNull();
	});

	it("Should keep the active command when another one is deleted", async () => {
		// Arrange
		installInMemoryBackend({
			commands: [
				new CommandBuilder().withId("cmd-1").build(),
				new CommandBuilder().withId("cmd-2").build(),
			],
		});
		commandStore.setState({ activeCommandId: "cmd-2" });

		// Act
		await sut("cmd-1");

		// Assert
		expect(commandStore.getState().activeCommandId).toBe("cmd-2");
	});

	it("Should dispose the terminal and the buffered logs of the command", async () => {
		// Arrange
		installInMemoryBackend({
			commands: [new CommandBuilder().withId("cmd-1").build()],
		});
		attachCommandOutput("cmd-1", document.createElement("div"));
		appendCommandOutput("cmd-1", ["a log line"]);

		// Act
		await sut("cmd-1");

		// Assert
		expect(recording.terminals.get("cmd-1")?.disposed).toBe(true);
		expect(commandOutputTail("cmd-1")).toEqual([]);
	});

	it("Should notify the user that the command was deleted", async () => {
		// Arrange
		installInMemoryBackend({
			commands: [new CommandBuilder().withId("cmd-1").build()],
		});

		// Act
		await sut("cmd-1");

		// Assert
		expect(toastSuccess).toHaveBeenCalledWith("toast.command.deleteSuccess");
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should still report success when the refresh that follows fails", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commands: [new CommandBuilder().withId("cmd-1").build()],
		});
		backend.data.getCommands = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut("cmd-1");

		// Assert
		expect(succeeded).toBe(true);
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should keep the command and notify the user when the backend rejects the deletion", async () => {
		// Arrange
		const command = new CommandBuilder().withId("cmd-1").build();
		const backend = installInMemoryBackend({ commands: [command] });
		backend.data.removeCommand = async () => {
			throw new Error("boom");
		};
		commandStore.setState({ activeCommandId: "cmd-1" });

		// Act
		const succeeded = await sut("cmd-1");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith("toast.command.deleteFailed: boom");
		expect(commandStore.getState().commands).toEqual([command]);
		expect(commandStore.getState().activeCommandId).toBe("cmd-1");
	});
});
