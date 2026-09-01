import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

const xterm = vi.hoisted(() => {
	type KeyHandler = (event: KeyboardEvent) => boolean;

	class Terminal {
		static latest: Terminal | null = null;

		element: HTMLElement | null = null;
		options: Record<string, unknown> = {};
		selection = "";
		keyHandler: KeyHandler = () => true;

		constructor() {
			Terminal.latest = this;
		}

		open(element: HTMLElement) {
			this.element = document.createElement("div");
			element.appendChild(this.element);
		}

		attachCustomKeyEventHandler(handler: KeyHandler) {
			this.keyHandler = handler;
		}

		getSelection() {
			return this.selection;
		}

		hasSelection() {
			return this.selection.length > 0;
		}

		loadAddon() {}
		writeln() {}
		reset() {}
		dispose() {}
	}

	return { Terminal };
});

vi.mock("@xterm/xterm", () => ({ Terminal: xterm.Terminal }));

vi.mock("@xterm/addon-fit", () => ({
	FitAddon: class {
		fit() {}
		dispose() {}
	},
}));

vi.mock("@xterm/addon-search", () => ({
	SearchAddon: class {
		findNext() {}
		findPrevious() {}
		clearDecorations() {}
		clearActiveDecoration() {}
		onDidChangeResults() {
			return { dispose() {} };
		}
		dispose() {}
	},
}));

vi.mock("@xterm/addon-web-links", () => ({
	WebLinksAddon: class {
		dispose() {}
	},
}));

import { xtermTerminal } from "@/commandOutput/adapters/xterm.ts";
import type { InMemoryBackend } from "@/contracts/adapters/inMemory.ts";
import { resetBackendServices } from "@/contracts/service.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";

const keyboardEvent = (overrides: Partial<KeyboardEvent> = {}): KeyboardEvent =>
	({
		altKey: false,
		ctrlKey: true,
		key: "c",
		repeat: false,
		type: "keydown",
		preventDefault: vi.fn(),
		...overrides,
	}) as unknown as KeyboardEvent;

describe("the xterm adapter clipboard behavior", () => {
	let backend: InMemoryBackend;

	beforeEach(() => {
		xterm.Terminal.latest = null;
		backend = installInMemoryBackend();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should copy the current selection once for Ctrl+C", () => {
		// Arrange
		const sut = xtermTerminal("cmd-1", "dark");
		sut.attach(document.createElement("div"));
		const terminal = xterm.Terminal.latest;
		if (!terminal) {
			throw new Error("Terminal was not created");
		}
		terminal.selection = "selected output";
		const event = keyboardEvent();

		// Act
		const handled = terminal.keyHandler(event);

		// Assert
		expect(handled).toBe(false);
		expect(event.preventDefault).toHaveBeenCalledOnce();
		expect(backend.state.clipboardText).toBe("selected output");
	});

	it("Should consume repeated Ctrl+C events without copying again", () => {
		// Arrange
		const sut = xtermTerminal("cmd-1", "dark");
		sut.attach(document.createElement("div"));
		const terminal = xterm.Terminal.latest;
		if (!terminal) {
			throw new Error("Terminal was not created");
		}
		terminal.selection = "selected output";
		const event = keyboardEvent({ repeat: true });

		// Act
		const handled = terminal.keyHandler(event);

		// Assert
		expect(handled).toBe(false);
		expect(backend.state.clipboardText).toBe("");
	});

	it("Should report whether the terminal has a selection", () => {
		// Arrange
		const sut = xtermTerminal("cmd-1", "dark").attach(
			document.createElement("div"),
		);
		const terminal = xterm.Terminal.latest;
		if (!terminal) {
			throw new Error("Terminal was not created");
		}

		// Act
		const withoutSelection = sut.hasSelection();
		terminal.selection = "selected output";
		const withSelection = sut.hasSelection();

		// Assert
		expect(withoutSelection).toBe(false);
		expect(withSelection).toBe(true);
	});

	it("Should restore normal key handling when detached", () => {
		// Arrange
		const sut = xtermTerminal("cmd-1", "dark");
		const attached = sut.attach(document.createElement("div"));
		const terminal = xterm.Terminal.latest;
		if (!terminal) {
			throw new Error("Terminal was not created");
		}
		terminal.selection = "selected output";

		// Act
		attached.detach();
		const handled = terminal.keyHandler(keyboardEvent());

		// Assert
		expect(handled).toBe(true);
		expect(backend.state.clipboardText).toBe("");
	});
});
