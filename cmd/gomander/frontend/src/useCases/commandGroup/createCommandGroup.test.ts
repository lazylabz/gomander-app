import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { createCommandGroup } from "@/useCases/commandGroup/createCommandGroup.ts";

describe("createCommandGroup", () => {
	const sut = createCommandGroup;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	const groupToCreate = {
		id: "group-1",
		projectId: "project-1",
		name: "A group",
		commandIds: ["cmd-2", "cmd-1"],
	};

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		commandGroupStore.setState({ commandGroups: [] });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should create the group naming the commands it was given, in that order, and refresh the groups", async () => {
		// Arrange
		installInMemoryBackend();

		// Act
		const succeeded = await sut(groupToCreate);

		// Assert
		expect(succeeded).toBe(true);
		const [group] = commandGroupStore.getState().commandGroups;
		expect(group.id).toBe("group-1");
		expect(group.name).toBe("A group");
		expect(group.commandIds).toEqual(["cmd-2", "cmd-1"]);
		expect(toastSuccess).toHaveBeenCalledWith(
			"toast.commandGroup.createSuccess",
		);
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
