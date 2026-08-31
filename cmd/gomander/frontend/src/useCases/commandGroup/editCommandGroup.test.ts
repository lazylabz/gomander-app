import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { editCommandGroup } from "@/useCases/commandGroup/editCommandGroup.ts";

describe("editCommandGroup", () => {
	const sut = editCommandGroup;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	const firstCommand = new CommandBuilder().withId("cmd-1").build();
	const secondCommand = new CommandBuilder().withId("cmd-2").build();

	const existingGroup = new CommandGroupBuilder()
		.withId("group-1")
		.withName("Old name")
		.withCommands(firstCommand)
		.build();

	const edit = {
		...existingGroup,
		name: "New name",
		commandIds: [secondCommand.id, firstCommand.id],
	};

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		commandGroupStore.setState({ commandGroups: [existingGroup] });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should save the edited group, naming the commands in the order it was given, and refresh the groups", async () => {
		// Arrange
		installInMemoryBackend({ commandGroups: [existingGroup] });

		// Act
		const succeeded = await sut(edit);

		// Assert
		expect(succeeded).toBe(true);
		const [group] = commandGroupStore.getState().commandGroups;
		expect(group.name).toBe("New name");
		expect(group.commandIds).toEqual(["cmd-2", "cmd-1"]);
		expect(toastSuccess).toHaveBeenCalledWith(
			"toast.commandGroup.updateSuccess",
		);
	});

	it("Should notify the user when the backend rejects the edit", async () => {
		// Arrange
		const backend = installInMemoryBackend({ commandGroups: [existingGroup] });
		backend.data.editCommandGroup = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut(edit);

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith(
			"toast.commandGroup.updateFailed: boom",
		);
		expect(commandGroupStore.getState().commandGroups).toEqual([existingGroup]);
	});
});
