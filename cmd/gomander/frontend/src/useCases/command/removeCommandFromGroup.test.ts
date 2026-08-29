import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import type { Command } from "@/contracts/types.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { removeCommandFromGroup } from "@/useCases/command/removeCommandFromGroup.ts";

describe("removeCommandFromGroup", () => {
	const sut = removeCommandFromGroup;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	const firstCommand = new CommandBuilder().withId("cmd-1").build();
	const secondCommand = new CommandBuilder().withId("cmd-2").build();

	const groupWith = (...commands: Command[]) =>
		new CommandGroupBuilder()
			.withId("group-1")
			.withCommands(...commands)
			.build();

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		commandGroupStore.setState({ commandGroups: [] });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should remove the command from the group and refresh the groups", async () => {
		// Arrange
		const group = groupWith(firstCommand, secondCommand);
		installInMemoryBackend({ commandGroups: [group] });
		commandGroupStore.setState({ commandGroups: [group] });

		// Act
		const succeeded = await sut("cmd-1", "group-1");

		// Assert
		expect(succeeded).toBe(true);
		expect(
			commandGroupStore.getState().commandGroups[0].commands.map((c) => c.id),
		).toEqual(["cmd-2"]);
		expect(toastSuccess).toHaveBeenCalledWith(
			"toast.command.removeFromGroupSuccess",
		);
	});

	it("Should refuse to remove the last command of the group", async () => {
		// Arrange
		const group = groupWith(firstCommand);
		const backend = installInMemoryBackend({ commandGroups: [group] });
		commandGroupStore.setState({ commandGroups: [group] });

		// Act
		const succeeded = await sut("cmd-1", "group-1");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith(
			"toast.commandGroup.cannotRemoveLast",
		);
		expect(backend.state.commandGroups[0].commands).toHaveLength(1);
	});

	it("Should report an unknown group", async () => {
		// Arrange
		installInMemoryBackend();

		// Act
		const succeeded = await sut("cmd-1", "group-1");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith("toast.commandGroup.notFound");
	});

	it("Should notify the user when the backend rejects the removal", async () => {
		// Arrange
		const group = groupWith(firstCommand, secondCommand);
		const backend = installInMemoryBackend({ commandGroups: [group] });
		commandGroupStore.setState({ commandGroups: [group] });
		backend.data.removeCommandFromGroup = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut("cmd-1", "group-1");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith(
			"toast.command.removeFromGroupFailed: boom",
		);
		expect(
			commandGroupStore.getState().commandGroups[0].commands.map((c) => c.id),
		).toEqual(["cmd-1", "cmd-2"]);
	});
});
