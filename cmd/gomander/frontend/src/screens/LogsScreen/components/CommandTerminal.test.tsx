import { act, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import type {
	RecordedTerminal,
	RecordingTerminals,
} from "@/commandOutput/adapters/recording.ts";
import { CommandTerminal } from "@/screens/LogsScreen/components/CommandTerminal.tsx";
import { installTranslations, withTranslation } from "@/testing/i18n.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import {
	installRecordingTerminals,
	resetTerminals,
} from "@/testing/terminals.ts";

describe("CommandTerminal", () => {
	let recording: RecordingTerminals;

	const terminalOf = (commandId: string): RecordedTerminal => {
		const terminal = recording.terminals.get(commandId);
		if (!terminal) {
			throw new Error(`No terminal was created for ${commandId}`);
		}
		return terminal;
	};

	const renderSut = () =>
		renderWithProviders(<CommandTerminal commandId="cmd-1" />);

	const openSearch = async (user: ReturnType<typeof renderSut>["user"]) => {
		await user.keyboard("{Meta>}f{/Meta}");
		return screen.getByRole("textbox");
	};

	beforeEach(async () => {
		await installTranslations();
		recording = installRecordingTerminals();
	});

	afterEach(() => {
		resetTerminals();
	});

	it("Should attach the terminal of the command when mounted", () => {
		// Arrange / Act
		const { container } = renderSut();

		// Assert
		const attachedTo = terminalOf("cmd-1").attachedTo;
		expect(attachedTo).not.toBeNull();
		expect(container.contains(attachedTo)).toBe(true);
	});

	it("Should detach the terminal without disposing it when unmounted", () => {
		// Arrange
		const { unmount } = renderSut();

		// Act
		unmount();

		// Assert
		expect(terminalOf("cmd-1").attachedTo).toBeNull();
		expect(terminalOf("cmd-1").disposed).toBe(false);
	});

	it("Should keep the search bar hidden until the shortcut is pressed", () => {
		// Arrange / Act
		renderSut();

		// Assert
		expect(screen.queryByRole("textbox")).not.toBeInTheDocument();
	});

	it("Should open the search bar with its input focused on Mod-F", async () => {
		// Arrange
		const { user } = renderSut();

		// Act
		const input = await openSearch(user);

		// Assert
		expect(input).toHaveFocus();
	});

	it("Should search the terminal incrementally as the user types", async () => {
		// Arrange
		const { user } = renderSut();
		const input = await openSearch(user);

		// Act
		await user.type(input, "ab");

		// Assert
		expect(terminalOf("cmd-1").searches).toEqual([
			{ query: "a", direction: "next", incremental: true },
			{ query: "ab", direction: "next", incremental: true },
		]);
	});

	it("Should move to the next and the previous match with the buttons", async () => {
		// Arrange
		const { user } = renderSut();
		const input = await openSearch(user);
		await user.type(input, "ab");
		const typed = terminalOf("cmd-1").searches.length;

		// Act
		await user.click(screen.getByRole("button", { name: "logs.searchNext" }));
		await user.click(screen.getByRole("button", { name: "logs.searchPrev" }));

		// Assert
		expect(terminalOf("cmd-1").searches.slice(typed)).toEqual([
			{ query: "ab", direction: "next", incremental: false },
			{ query: "ab", direction: "previous", incremental: false },
		]);
	});

	it("Should show how many matches the terminal reports", async () => {
		// Arrange
		withTranslation("logs.matches", "{{count}} matches");
		const { user } = renderSut();
		const input = await openSearch(user);
		await user.type(input, "ab");

		// Act
		act(() => terminalOf("cmd-1").reportResults(3));

		// Assert
		expect(screen.getByText("3 matches")).toBeInTheDocument();
	});

	it("Should hide the search bar and clear the search when closed", async () => {
		// Arrange
		const { user } = renderSut();
		const input = await openSearch(user);
		await user.type(input, "ab");
		const clearedBefore = terminalOf("cmd-1").searchesCleared;

		// Act
		await user.click(screen.getByRole("button", { name: "common.close" }));

		// Assert
		expect(screen.queryByRole("textbox")).not.toBeInTheDocument();
		expect(terminalOf("cmd-1").searchesCleared).toBeGreaterThan(clearedBefore);
	});
});
