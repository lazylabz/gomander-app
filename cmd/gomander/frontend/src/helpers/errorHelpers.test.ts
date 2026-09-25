import { describe, expect, it } from "vitest";

import { parseError } from "@/helpers/errorHelpers.ts";

describe("parseError", () => {
	const sut = parseError;

	it.each([
		["an Error", new Error("boom")],
		["a string", "boom"],
	])("Should read the message of %s", (_, error) => {
		// Act
		const message = sut(error);

		// Assert
		expect(message).toBe("boom");
	});

	it.each([
		["an Error", new Error("boom")],
		["a string", "boom"],
	])("Should put the prefix before the message of %s", (_, error) => {
		// Act
		const message = sut(error, "Could not save");

		// Assert
		expect(message).toBe("Could not save: boom");
	});

	it("Should fall back to a generic message, without the prefix, when the error is neither an Error nor a string", () => {
		// Act
		const message = sut({ code: 1 }, "Could not save");

		// Assert
		expect(message).toBe("Unknown error");
	});
});
