import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { stopCommandGroup } from "@/useCases/commandGroup/stopCommandGroup.ts";

describe("stopCommandGroup", () => {
	const sut = stopCommandGroup;

	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should ask the backend to stop the group", async () => {
		// Arrange
		const backend = installInMemoryBackend({ runningGroupIds: ["group-1"] });

		// Act
		const succeeded = await sut("group-1");

		// Assert
		expect(succeeded).toBe(true);
		expect(backend.state.runningGroupIds).toEqual([]);
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should notify the user when the group cannot be stopped", async () => {
		// Arrange
		const backend = installInMemoryBackend({ runningGroupIds: ["group-1"] });
		backend.data.stopCommandGroup = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut("group-1");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith(
			"toast.commandGroup.stopFailed: boom",
		);
	});
});
