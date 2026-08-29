import { beforeEach, describe, expect, it, vi } from "vitest";

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
	beforeEach(() => {
		xterm.Terminal.latest = null;
	});

	it("Should copy the current selection once for Ctrl+C", () => {
		// Arrange
		const copyText = vi.fn();
		const output = xtermTerminal("cmd-1", "dark");
		output.attach(document.createElement("div"), copyText);
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
		expect(copyText).toHaveBeenCalledExactlyOnceWith("selected output");
	});

	it("Should consume repeated Ctrl+C events without copying again", () => {
		// Arrange
		const copyText = vi.fn();
		const output = xtermTerminal("cmd-1", "dark");
		output.attach(document.createElement("div"), copyText);
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
		expect(copyText).not.toHaveBeenCalled();
	});

	it("Should restore normal key handling when detached", () => {
		// Arrange
		const copyText = vi.fn();
		const output = xtermTerminal("cmd-1", "dark");
		const attached = output.attach(document.createElement("div"), copyText);
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
		expect(copyText).not.toHaveBeenCalled();
	});
});
