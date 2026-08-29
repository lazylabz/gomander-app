import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { duplicateCommand } from "@/useCases/command/duplicateCommand.ts";

describe("duplicateCommand", () => {
	const sut = duplicateCommand;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		commandStore.setState({
			commands: [],
			commandsStatus: {},
			activeCommandId: null,
		});
		commandGroupStore.setState({ commandGroups: [] });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should refresh the commands in the store with the duplicate", async () => {
		// Arrange
		installInMemoryBackend({
			commands: [new CommandBuilder().withId("cmd-1").build()],
		});

		// Act
		const succeeded = await sut("cmd-1");

		// Assert
		expect(succeeded).toBe(true);
		expect(commandStore.getState().commands).toHaveLength(2);
	});

	it("Should refresh the target group when the duplicate lands in one", async () => {
		// Arrange
		const command = new CommandBuilder().withId("cmd-1").build();
		installInMemoryBackend({
			commands: [command],
			commandGroups: [
				new CommandGroupBuilder()
					.withId("group-1")
					.withCommands(command)
					.build(),
			],
		});

		// Act
		await sut("cmd-1", "group-1");

		// Assert
		expect(commandGroupStore.getState().commandGroups[0].commands).toHaveLength(
			2,
		);
	});

	it("Should clear the active command and notify the user", async () => {
		// Arrange
		installInMemoryBackend({
			commands: [new CommandBuilder().withId("cmd-1").build()],
		});
		commandStore.setState({ activeCommandId: "cmd-1" });

		// Act
		await sut("cmd-1");

		// Assert
		expect(commandStore.getState().activeCommandId).toBeNull();
		expect(toastSuccess).toHaveBeenCalledWith("toast.command.duplicateSuccess");
	});

	it("Should notify the user when the backend rejects the duplication", async () => {
		// Arrange
		const command = new CommandBuilder().withId("cmd-1").build();
		const backend = installInMemoryBackend({ commands: [command] });
		backend.data.duplicateCommand = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut("cmd-1");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith(
			"toast.command.duplicateFailed: boom",
		);
		expect(toastSuccess).not.toHaveBeenCalled();
		expect(commandStore.getState().commands).toEqual([command]);
	});
});
