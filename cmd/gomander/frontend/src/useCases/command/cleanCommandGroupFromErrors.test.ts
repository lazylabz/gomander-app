import { beforeEach, describe, expect, it } from "vitest";

import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { cleanCommandGroupFromErrors } from "@/useCases/command/cleanCommandGroupFromErrors.ts";

describe("cleanCommandGroupFromErrors", () => {
	const sut = cleanCommandGroupFromErrors;

	const insideGroup = new CommandBuilder().withId("cmd-1").build();
	const alsoInsideGroup = new CommandBuilder().withId("cmd-2").build();
	const outsideGroup = new CommandBuilder().withId("cmd-3").build();

	beforeEach(() => {
		commandGroupStore.setState({
			commandGroups: [
				new CommandGroupBuilder()
					.withId("group-1")
					.withCommands(insideGroup, alsoInsideGroup)
					.build(),
			],
		});
		commandStore.setState({
			commandIdsWithErrors: ["cmd-1", "cmd-2", "cmd-3"],
		});
	});

	it("Should clear the errors of the group's own commands", () => {
		// Act
		sut("group-1");

		// Assert
		expect(commandStore.getState().commandIdsWithErrors).toEqual(["cmd-3"]);
	});

	it("Should leave the errors of a command outside the group alone", () => {
		// Arrange
		commandStore.setState({ commandIdsWithErrors: [outsideGroup.id] });

		// Act
		sut("group-1");

		// Assert
		expect(commandStore.getState().commandIdsWithErrors).toEqual(["cmd-3"]);
	});

	it("Should leave every error alone when the group is unknown", () => {
		// Act
		sut("group-that-does-not-exist");

		// Assert
		expect(commandStore.getState().commandIdsWithErrors).toEqual([
			"cmd-1",
			"cmd-2",
			"cmd-3",
		]);
	});
});
