import { describe, expect, it } from "vitest";

import { resolveCommands } from "@/helpers/commandHelpers.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";

describe("resolveCommands", () => {
	it("should follow the order the ids name, not the order the commands are loaded in", () => {
		// Arrange
		const first = new CommandBuilder().withId("cmd-1").build();
		const second = new CommandBuilder().withId("cmd-2").build();

		// Act
		const resolved = resolveCommands([second.id, first.id], [first, second]);

		// Assert
		expect(resolved).toEqual([second, first]);
	});

	it("should drop an id no loaded command answers to", () => {
		// Arrange
		const command = new CommandBuilder().withId("cmd-1").build();

		// Act
		const resolved = resolveCommands(["unknown", command.id], [command]);

		// Assert
		expect(resolved).toEqual([command]);
	});
});
