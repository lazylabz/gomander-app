import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { runCommandGroup } from "@/useCases/commandGroup/runCommandGroup.ts";

describe("runCommandGroup", () => {
	const sut = runCommandGroup;

	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		commandGroupStore.setState({
			commandGroups: [new CommandGroupBuilder().withId("group-1").build()],
		});
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should ask the backend to run the group", async () => {
		// Arrange
		const backend = installInMemoryBackend();

		// Act
		const succeeded = await sut("group-1");

		// Assert
		expect(succeeded).toBe(true);
		expect(backend.state.runningGroupIds).toEqual(["group-1"]);
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should notify the user when the group cannot be run", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.runCommandGroup = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut("group-1");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith(
			"toast.commandGroup.runFailed: boom",
		);
	});
});
