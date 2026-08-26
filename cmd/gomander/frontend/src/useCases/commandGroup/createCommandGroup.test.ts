import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { createCommandGroup } from "@/useCases/commandGroup/createCommandGroup.ts";

describe("createCommandGroup", () => {
	const sut = createCommandGroup;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	const firstCommand = new CommandBuilder().withId("cmd-1").build();
	const secondCommand = new CommandBuilder().withId("cmd-2").build();

	const groupToCreate = {
		id: "group-1",
		projectId: "project-1",
		name: "A group",
		commands: ["cmd-1", "cmd-2"],
	};

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		commandStore.setState({ commands: [firstCommand, secondCommand] });
		commandGroupStore.setState({ commandGroups: [] });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should create the group with the commands it was given and refresh the groups", async () => {
		// Arrange
		installInMemoryBackend();

		// Act
		const succeeded = await sut(groupToCreate);

		// Assert
		expect(succeeded).toBe(true);
		const [group] = commandGroupStore.getState().commandGroups;
		expect(group.id).toBe("group-1");
		expect(group.name).toBe("A group");
		expect(group.commands).toEqual([firstCommand, secondCommand]);
		expect(toastSuccess).toHaveBeenCalledWith(
			"toast.commandGroup.createSuccess",
		);
	});

	it("Should leave out a command it does not know about", async () => {
		// Arrange
		const backend = installInMemoryBackend();

		// Act
		await sut({ ...groupToCreate, commands: ["cmd-1", "unknown"] });

		// Assert
		expect(backend.state.commandGroups[0].commands.map((c) => c.id)).toEqual([
			"cmd-1",
		]);
	});

	it("Should notify the user when the backend rejects the creation", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.createCommandGroup = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut(groupToCreate);

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith(
			"toast.commandGroup.createFailed: boom",
		);
		expect(commandGroupStore.getState().commandGroups).toEqual([]);
	});
});
