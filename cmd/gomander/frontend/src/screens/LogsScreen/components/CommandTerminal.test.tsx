import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { I18nextProvider } from "react-i18next";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { CommandTerminal } from "@/screens/LogsScreen/components/CommandTerminal.tsx";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import {
	installRecordingTerminals,
	resetTerminals,
} from "@/testing/terminals.ts";

class TestResizeObserver {
	disconnect() {}
	observe() {}
	unobserve() {}
}

describe("CommandTerminal", () => {
	let root: Root;
	let host: HTMLDivElement;

	beforeEach(async () => {
		Object.assign(globalThis, {
			IS_REACT_ACT_ENVIRONMENT: true,
			ResizeObserver: TestResizeObserver,
		});
		await installTranslations();

		host = document.createElement("div");
		document.body.appendChild(host);
		root = createRoot(host);
	});

	afterEach(() => {
		act(() => root.unmount());
		host.remove();
		resetTerminals();
		resetBackendServices();
	});

	it("Should copy the selected output from the context menu", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		const recording = installRecordingTerminals();
		const sut = <CommandTerminal commandId="cmd-1" />;
		await act(async () =>
			root.render(<I18nextProvider i18n={i18n}>{sut}</I18nextProvider>),
		);

		const terminal = recording.terminals.get("cmd-1");
		if (!terminal) {
			throw new Error("Terminal was not created");
		}
		terminal.selection = "selected output";

		const trigger = host.querySelector<HTMLElement>(
			'[data-slot="context-menu-trigger"]',
		);
		if (!trigger) {
			throw new Error("Context menu trigger was not rendered");
		}

		// Act
		await act(async () => {
			trigger.dispatchEvent(
				new MouseEvent("contextmenu", {
					bubbles: true,
					clientX: 10,
					clientY: 10,
				}),
			);
		});

		const copyItem = document.body.querySelector<HTMLElement>(
			'[data-slot="context-menu-item"]',
		);
		if (!copyItem) {
			throw new Error("Copy menu item was not rendered");
		}

		await act(async () => copyItem.click());

		// Assert
		expect(copyItem.hasAttribute("data-disabled")).toBe(false);
		expect(backend.state.clipboardText).toBe("selected output");
	});
});
