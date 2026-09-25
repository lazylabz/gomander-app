import { describe, expect, it } from "vitest";

import { matchesShortcut } from "@/hooks/useShortcut.ts";

describe("matchesShortcut", () => {
	const sut = matchesShortcut;

	it("Should match a bare key regardless of its case", () => {
		// Arrange
		const event = new KeyboardEvent("keydown", { key: "k" });

		// Act
		const matches = sut("K", event);

		// Assert
		expect(matches).toBe(true);
	});

	it("Should not match a different key", () => {
		// Arrange
		const event = new KeyboardEvent("keydown", { key: "j", ctrlKey: true });

		// Act
		const matches = sut("Control-K", event);

		// Assert
		expect(matches).toBe(false);
	});

	it.each([
		["Control-K", { ctrlKey: true }],
		["Shift-K", { shiftKey: true }],
		["Alt-K", { altKey: true }],
		["Mod-K", { metaKey: true }],
		["Mod-K", { ctrlKey: true }],
	] as const)("Should match %s when pressed with %o", (shortcut, modifiers) => {
		// Arrange
		const event = new KeyboardEvent("keydown", { key: "k", ...modifiers });

		// Act
		const matches = sut(shortcut, event);

		// Assert
		expect(matches).toBe(true);
	});

	it.each([
		["Control-K", {}],
		["Control-K", { metaKey: true }],
		["Shift-K", { ctrlKey: true }],
		["Alt-K", { shiftKey: true }],
		["Mod-K", { altKey: true }],
	] as const)("Should not match %s when pressed with %o", (shortcut, modifiers) => {
		// Arrange
		const event = new KeyboardEvent("keydown", { key: "k", ...modifiers });

		// Act
		const matches = sut(shortcut, event);

		// Assert
		expect(matches).toBe(false);
	});
});
