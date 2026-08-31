import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { editCommand } from "@/useCases/command/editCommand.ts";

describe("editCommand", () => {
	const sut = editCommand;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		commandStore.setState({ commands: [], commandsStatus: {} });
		commandGroupStore.setState({ commandGroups: [] });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should persist the edited command and refresh the commands in the store", async () => {
		// Arrange
		const command = new CommandBuilder()
			.withId("cmd-1")
			.withName("old")
			.build();
		installInMemoryBackend({ commands: [command] });
		const edited = { ...command, name: "new" };

		// Act
		const succeeded = await sut(edited);

		// Assert
		expect(succeeded).toBe(true);
		expect(commandStore.getState().commands).toEqual([edited]);
	});

	it("Should notify the user that the command was updated", async () => {
		// Arrange
		const command = new CommandBuilder().withId("cmd-1").build();
		installInMemoryBackend({ commands: [command] });

		// Act
		await sut(command);

		// Assert
		expect(toastSuccess).toHaveBeenCalledWith("toast.command.updateSuccess");
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should notify the user when the backend rejects the edit", async () => {
		// Arrange
		const command = new CommandBuilder()
			.withId("cmd-1")
			.withName("old")
			.build();
		const backend = installInMemoryBackend({ commands: [command] });
		backend.data.editCommand = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut({ ...command, name: "new" });

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith("toast.command.updateFailed: boom");
		expect(commandStore.getState().commands).toEqual([command]);
	});
});
