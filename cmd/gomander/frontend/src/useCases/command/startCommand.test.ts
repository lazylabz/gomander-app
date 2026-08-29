import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { commandStore } from "@/store/commandStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { startCommand } from "@/useCases/command/startCommand.ts";

describe("startCommand", () => {
	const sut = startCommand;

	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		commandStore.setState({ commandIdsWithErrors: [] });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should ask the backend to run the command", async () => {
		// Arrange
		const backend = installInMemoryBackend();

		// Act
		const succeeded = await sut("cmd-1");

		// Assert
		expect(succeeded).toBe(true);
		expect(backend.state.runningCommandIds).toEqual(["cmd-1"]);
	});

	it("Should clear the error the command was left in", async () => {
		// Arrange
		installInMemoryBackend();
		commandStore.setState({ commandIdsWithErrors: ["cmd-1", "cmd-2"] });

		// Act
		await sut("cmd-1");

		// Assert
		expect(commandStore.getState().commandIdsWithErrors).toEqual(["cmd-2"]);
	});

	it("Should notify the user when the command cannot be run", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.runCommand = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut("cmd-1");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith("toast.command.runFailed: boom");
	});
});
