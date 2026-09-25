import { describe, expect, it } from "vitest";

import {
	addCommand,
	moveAcrossContainers,
	removeCommand,
	reorderAdded,
} from "@/components/modals/CommandGroup/common/CommandGroupCommandsField/commandGroupSelection.ts";
import {
	ADDED_COMMANDS,
	AVAILABLE_COMMANDS,
} from "@/components/modals/CommandGroup/common/CommandGroupCommandsField/constants.ts";

const allCommandIds = ["cmd-1", "cmd-2", "cmd-3", "cmd-4"];

describe("addCommand", () => {
	const sut = addCommand;

	it("Should append the command after the ones already selected", () => {
		// Act
		const selected = sut(["cmd-2"], "cmd-1");

		// Assert
		expect(selected).toEqual(["cmd-2", "cmd-1"]);
	});
});

describe("removeCommand", () => {
	const sut = removeCommand;

	it("Should drop the command and keep the order of the rest", () => {
		// Act
		const selected = sut(["cmd-1", "cmd-2", "cmd-3"], "cmd-2");

		// Assert
		expect(selected).toEqual(["cmd-1", "cmd-3"]);
	});
});

describe("moveAcrossContainers", () => {
	const sut = moveAcrossContainers;

	it("Should add an available command dropped on the added container", () => {
		// Arrange
		const selectedIds = ["cmd-1"];

		// Act
		const selected = sut(selectedIds, allCommandIds, {
			activeId: "cmd-2",
			overId: ADDED_COMMANDS,
		});

		// Assert
		expect(selected).toEqual(["cmd-1", "cmd-2"]);
	});

	it("Should add an available command dropped on an added command", () => {
		// Arrange
		const selectedIds = ["cmd-1"];

		// Act
		const selected = sut(selectedIds, allCommandIds, {
			activeId: "cmd-2",
			overId: "cmd-1",
		});

		// Assert
		expect(selected).toEqual(["cmd-1", "cmd-2"]);
	});

	it("Should remove an added command dropped on the available container", () => {
		// Arrange
		const selectedIds = ["cmd-1", "cmd-2"];

		// Act
		const selected = sut(selectedIds, allCommandIds, {
			activeId: "cmd-1",
			overId: AVAILABLE_COMMANDS,
		});

		// Assert
		expect(selected).toEqual(["cmd-2"]);
	});

	it("Should remove an added command dropped on an available command", () => {
		// Arrange
		const selectedIds = ["cmd-1", "cmd-2"];

		// Act
		const selected = sut(selectedIds, allCommandIds, {
			activeId: "cmd-1",
			overId: "cmd-3",
		});

		// Assert
		expect(selected).toEqual(["cmd-2"]);
	});

	it("Should hand back the same array when the drag stays within one container", () => {
		// Arrange
		const selectedIds = ["cmd-1", "cmd-2"];

		// Act
		const selected = sut(selectedIds, allCommandIds, {
			activeId: "cmd-1",
			overId: "cmd-2",
		});

		// Assert
		expect(selected).toBe(selectedIds);
	});

	it("Should ignore a dragged id that is not a command", () => {
		// Arrange
		const selectedIds = ["cmd-1"];

		// Act
		const selected = sut(selectedIds, allCommandIds, {
			activeId: "unknown",
			overId: ADDED_COMMANDS,
		});

		// Assert
		expect(selected).toBe(selectedIds);
	});

	it("Should keep an added command dropped over an unknown id", () => {
		// Arrange
		const selectedIds = ["cmd-1"];

		// Act
		const selected = sut(selectedIds, allCommandIds, {
			activeId: "cmd-1",
			overId: "unknown",
		});

		// Assert
		expect(selected).toBe(selectedIds);
	});
});

describe("reorderAdded", () => {
	const sut = reorderAdded;

	it("Should move an added command to the position of the one it is dropped on", () => {
		// Arrange
		const selectedIds = ["cmd-1", "cmd-2", "cmd-3"];

		// Act
		const selected = sut(selectedIds, {
			activeId: "cmd-3",
			overId: "cmd-1",
		});

		// Assert
		expect(selected).toEqual(["cmd-3", "cmd-1", "cmd-2"]);
	});

	it("Should hand back the same array when the drag is not between added commands", () => {
		// Arrange
		const selectedIds = ["cmd-1", "cmd-2"];

		// Act
		const selected = sut(selectedIds, {
			activeId: "cmd-3",
			overId: "cmd-1",
		});

		// Assert
		expect(selected).toBe(selectedIds);
	});

	it("Should hand back the same array when dropped on the added container itself", () => {
		// Arrange
		const selectedIds = ["cmd-1", "cmd-2"];

		// Act
		const selected = sut(selectedIds, {
			activeId: "cmd-1",
			overId: ADDED_COMMANDS,
		});

		// Assert
		expect(selected).toBe(selectedIds);
	});

	it("Should hand back the same array when dropped on itself", () => {
		// Arrange
		const selectedIds = ["cmd-1", "cmd-2"];

		// Act
		const selected = sut(selectedIds, {
			activeId: "cmd-2",
			overId: "cmd-2",
		});

		// Assert
		expect(selected).toBe(selectedIds);
	});
});
