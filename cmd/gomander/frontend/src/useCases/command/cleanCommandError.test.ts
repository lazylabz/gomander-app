import { beforeEach, describe, expect, it } from "vitest";

import { commandStore } from "@/store/commandStore.ts";
import { resetStores } from "@/testing/stores.ts";
import { cleanCommandError } from "@/useCases/command/cleanCommandError.ts";

describe("cleanCommandError", () => {
	const sut = cleanCommandError;

	beforeEach(() => {
		resetStores();
	});

	it("Should clear the error of the command and keep the others", () => {
		// Arrange
		commandStore.setState({ commandIdsWithErrors: ["cmd-1", "cmd-2"] });

		// Act
		sut("cmd-1");

		// Assert
		expect(commandStore.getState().commandIdsWithErrors).toEqual(["cmd-2"]);
	});

	it("Should leave the errors alone when the command has none", () => {
		// Arrange
		const commandIdsWithErrors = ["cmd-2"];
		commandStore.setState({ commandIdsWithErrors });

		// Act
		sut("cmd-1");

		// Assert
		expect(commandStore.getState().commandIdsWithErrors).toBe(
			commandIdsWithErrors,
		);
	});
});
