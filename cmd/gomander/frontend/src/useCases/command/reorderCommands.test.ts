import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { commandStore } from "@/store/commandStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { reorderCommands } from "@/useCases/command/reorderCommands.ts";

describe("reorderCommands", () => {
	const sut = reorderCommands;

	const toastError = vi.spyOn(toast, "error");

	const firstCommand = new CommandBuilder().withId("cmd-1").build();
	const secondCommand = new CommandBuilder().withId("cmd-2").build();

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		commandStore.setState({
			commands: [firstCommand, secondCommand],
			commandsStatus: {},
		});
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should persist the new order and refresh the commands in the store", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commands: [firstCommand, secondCommand],
		});

		// Act
		const succeeded = await sut(["cmd-2", "cmd-1"]);

		// Assert
		expect(succeeded).toBe(true);
		expect(backend.state.commands.map((c) => c.id)).toEqual(["cmd-2", "cmd-1"]);
		expect(commandStore.getState().commands.map((c) => c.id)).toEqual([
			"cmd-2",
			"cmd-1",
		]);
	});

	it("Should apply the new order before the backend answers", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commands: [firstCommand, secondCommand],
		});
		let confirmReorder = () => {};
		backend.data.reorderCommands = () =>
			new Promise((resolve) => {
				confirmReorder = resolve;
			});

		// Act
		const pending = sut(["cmd-2", "cmd-1"]);

		// Assert
		expect(commandStore.getState().commands.map((c) => c.id)).toEqual([
			"cmd-2",
			"cmd-1",
		]);
		confirmReorder();
		await pending;
	});

	it("Should restore the persisted order and notify the user when the backend rejects it", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commands: [firstCommand, secondCommand],
		});
		backend.data.reorderCommands = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut(["cmd-2", "cmd-1"]);

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith(
			"toast.command.reorderFailed: boom",
		);
		expect(commandStore.getState().commands.map((c) => c.id)).toEqual([
			"cmd-1",
			"cmd-2",
		]);
	});
});
