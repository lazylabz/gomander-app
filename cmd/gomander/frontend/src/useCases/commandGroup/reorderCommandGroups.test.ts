import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { reorderCommandGroups } from "@/useCases/commandGroup/reorderCommandGroups.ts";

describe("reorderCommandGroups", () => {
	const sut = reorderCommandGroups;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	const firstGroup = new CommandGroupBuilder().withId("group-1").build();
	const secondGroup = new CommandGroupBuilder().withId("group-2").build();

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		commandGroupStore.setState({ commandGroups: [firstGroup, secondGroup] });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should persist the new order and refresh the groups in the store", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commandGroups: [firstGroup, secondGroup],
		});

		// Act
		const succeeded = await sut(["group-2", "group-1"]);

		// Assert
		expect(succeeded).toBe(true);
		expect(backend.state.commandGroups.map((g) => g.id)).toEqual([
			"group-2",
			"group-1",
		]);
		expect(commandGroupStore.getState().commandGroups.map((g) => g.id)).toEqual(
			["group-2", "group-1"],
		);
		expect(toastSuccess).toHaveBeenCalledWith(
			"toast.commandGroup.reorderSuccess",
		);
	});

	it("Should apply the new order before the backend answers", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commandGroups: [firstGroup, secondGroup],
		});
		let confirmReorder = () => {};
		backend.data.reorderCommandGroups = () =>
			new Promise((resolve) => {
				confirmReorder = resolve;
			});

		// Act
		const pending = sut(["group-2", "group-1"]);

		// Assert
		expect(commandGroupStore.getState().commandGroups.map((g) => g.id)).toEqual(
			["group-2", "group-1"],
		);
		confirmReorder();
		await pending;
	});

	it("Should restore the persisted order and notify the user when the backend rejects it", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commandGroups: [firstGroup, secondGroup],
		});
		backend.data.reorderCommandGroups = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut(["group-2", "group-1"]);

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith(
			"toast.commandGroup.reorderFailed: boom",
		);
		expect(commandGroupStore.getState().commandGroups.map((g) => g.id)).toEqual(
			["group-1", "group-2"],
		);
	});
});
