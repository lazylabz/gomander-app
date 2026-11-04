import { describe, expect, it } from "vitest";

import { zoomShortcutOf } from "@/helpers/zoomShortcut.ts";

describe("zoomShortcutOf", () => {
	const sut = zoomShortcutOf;

	it.each([
		["=", "in"],
		["+", "in"],
		["-", "out"],
		["0", "reset"],
	])("Should read Ctrl and %s as zooming %s", (key, direction) => {
		// Act
		const result = sut(new KeyboardEvent("keydown", { key, ctrlKey: true }));

		// Assert
		expect(result).toBe(direction);
	});

	it.each([
		["without Ctrl", { key: "=" }],
		["with Alt too", { key: "=", ctrlKey: true, altKey: true }],
		["with Meta too", { key: "=", ctrlKey: true, metaKey: true }],
		["on another key", { key: "k", ctrlKey: true }],
	])("Should ignore a press %s", (_, init) => {
		// Act
		const result = sut(new KeyboardEvent("keydown", init));

		// Assert
		expect(result).toBeNull();
	});
});
