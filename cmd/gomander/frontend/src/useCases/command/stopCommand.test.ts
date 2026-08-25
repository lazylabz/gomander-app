import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { stopCommand } from "@/useCases/command/stopCommand.ts";

describe("stopCommand", () => {
	const sut = stopCommand;

	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should ask the backend to stop the command", async () => {
		// Arrange
		const backend = installInMemoryBackend({ runningCommandIds: ["cmd-1"] });

		// Act
		const succeeded = await sut("cmd-1");

		// Assert
		expect(succeeded).toBe(true);
		expect(backend.state.runningCommandIds).toEqual([]);
	});

	it("Should notify the user when the command cannot be stopped", async () => {
		// Arrange
		const backend = installInMemoryBackend({ runningCommandIds: ["cmd-1"] });
		backend.data.stopCommand = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut("cmd-1");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith("toast.command.stopFailed: boom");
	});
});
