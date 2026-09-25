import { describe, expect, it } from "vitest";
import {
	isGroupSelectable,
	keepSelectableGroups,
	narrowBlueprint,
} from "@/components/modals/Project/common/blueprintSelection.ts";
import type {
	CommandBlueprint,
	CommandGroupBlueprint,
} from "@/contracts/types.ts";
import { ProjectBlueprintBuilder } from "@/testing/builders/projectBlueprint.ts";

const aCommand = (id: string): CommandBlueprint => ({
	id,
	name: id,
	command: `echo ${id}`,
	workingDirectory: "/app",
});

const aGroup = (id: string, commandIds: string[]): CommandGroupBlueprint => ({
	id,
	name: id,
	commandIds,
});

describe("isGroupSelectable", () => {
	const sut = isGroupSelectable;

	it("Should allow a group when at least one of its commands is selected", () => {
		// Arrange
		const group = aGroup("group-1", ["cmd-1", "cmd-2"]);

		// Act
		const selectable = sut(group, ["cmd-2"]);

		// Assert
		expect(selectable).toBe(true);
	});

	it("Should reject a group when none of its commands is selected", () => {
		// Arrange
		const group = aGroup("group-1", ["cmd-1", "cmd-2"]);

		// Act
		const selectable = sut(group, ["cmd-3"]);

		// Assert
		expect(selectable).toBe(false);
	});
});

describe("keepSelectableGroups", () => {
	const sut = keepSelectableGroups;

	it("Should drop the groups left without a selected command and keep the order of the rest", () => {
		// Arrange
		const blueprint = new ProjectBlueprintBuilder()
			.withCommandGroups(
				aGroup("group-1", ["cmd-1"]),
				aGroup("group-2", ["cmd-2"]),
				aGroup("group-3", ["cmd-3"]),
			)
			.build();

		// Act
		const groupIds = sut(
			blueprint,
			["group-3", "group-2", "group-1"],
			["cmd-1", "cmd-3"],
		);

		// Assert
		expect(groupIds).toEqual(["group-3", "group-1"]);
	});

	it("Should drop a group id the blueprint does not contain", () => {
		// Arrange
		const blueprint = new ProjectBlueprintBuilder()
			.withCommandGroups(aGroup("group-1", ["cmd-1"]))
			.build();

		// Act
		const groupIds = sut(blueprint, ["unknown", "group-1"], ["cmd-1"]);

		// Assert
		expect(groupIds).toEqual(["group-1"]);
	});
});

describe("narrowBlueprint", () => {
	const sut = narrowBlueprint;

	it("Should keep only the selected commands and groups", () => {
		// Arrange
		const blueprint = new ProjectBlueprintBuilder()
			.withName("Imported")
			.withWorkingDirectory("/imported")
			.withCommands(aCommand("cmd-1"), aCommand("cmd-2"))
			.withCommandGroups(
				aGroup("group-1", ["cmd-1"]),
				aGroup("group-2", ["cmd-2"]),
			)
			.build();

		// Act
		const narrowed = sut(blueprint, ["cmd-2"], ["group-2"]);

		// Assert
		expect(narrowed).toEqual({
			name: "Imported",
			workingDirectory: "/imported",
			commands: [aCommand("cmd-2")],
			commandGroups: [aGroup("group-2", ["cmd-2"])],
		});
	});
});
