import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { commandStore } from "@/store/commandStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { createCommand } from "@/useCases/command/createCommand.ts";

describe("createCommand", () => {
	const sut = createCommand;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		commandStore.setState({ commands: [], commandsStatus: {} });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should persist the command and refresh the commands in the store", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		const command = new CommandBuilder().withId("cmd-1").build();

		// Act
		const succeeded = await sut(command);

		// Assert
		expect(succeeded).toBe(true);
		expect(backend.state.commands).toEqual([command]);
		expect(commandStore.getState().commands).toEqual([command]);
	});

	it("Should notify the user that the command was created", async () => {
		// Arrange
		installInMemoryBackend();

		// Act
		await sut(new CommandBuilder().build());

		// Assert
		expect(toastSuccess).toHaveBeenCalledWith("toast.command.createSuccess");
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should notify the user when the backend rejects the command", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.addCommand = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut(new CommandBuilder().build());

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith("toast.command.createFailed: boom");
		expect(toastSuccess).not.toHaveBeenCalled();
		expect(commandStore.getState().commands).toEqual([]);
	});
});
