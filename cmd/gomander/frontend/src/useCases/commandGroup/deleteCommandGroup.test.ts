import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { deleteCommandGroup } from "@/useCases/commandGroup/deleteCommandGroup.ts";

describe("deleteCommandGroup", () => {
	const sut = deleteCommandGroup;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	const command = new CommandBuilder().withId("cmd-1").build();
	const group = new CommandGroupBuilder()
		.withId("group-1")
		.withCommands(command)
		.build();

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		commandGroupStore.setState({ commandGroups: [group] });
		commandStore.setState({ commands: [command], commandsStatus: {} });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should delete the group and refresh the groups", async () => {
		// Arrange
		installInMemoryBackend({ commands: [command], commandGroups: [group] });

		// Act
		const succeeded = await sut("group-1");

		// Assert
		expect(succeeded).toBe(true);
		expect(commandGroupStore.getState().commandGroups).toEqual([]);
		expect(toastSuccess).toHaveBeenCalledWith(
			"toast.commandGroup.deleteSuccess",
		);
	});

	it("Should keep the commands the group held", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commands: [command],
			commandGroups: [group],
		});

		// Act
		await sut("group-1");

		// Assert
		expect(backend.state.commands).toEqual([command]);
	});

	it("Should notify the user when the backend rejects the deletion", async () => {
		// Arrange
		const backend = installInMemoryBackend({ commandGroups: [group] });
		backend.data.deleteCommandGroup = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut("group-1");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith(
			"toast.commandGroup.deleteFailed: boom",
		);
		expect(commandGroupStore.getState().commandGroups).toEqual([group]);
	});
});
