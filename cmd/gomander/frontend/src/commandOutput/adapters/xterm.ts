import "@xterm/xterm/css/xterm.css";
import { FitAddon } from "@xterm/addon-fit";
import type { ISearchOptions } from "@xterm/addon-search";
import { SearchAddon } from "@xterm/addon-search";
import { WebLinksAddon } from "@xterm/addon-web-links";
import { type ITheme, Terminal } from "@xterm/xterm";

import {
	type OutputTerminal,
	TERMINAL_SCROLLBACK,
	type TerminalFactory,
	type TerminalTheme,
} from "@/commandOutput/ports.ts";
import { externalBrowserService } from "@/contracts/service.ts";

const XTERM_THEMES: Record<TerminalTheme, ITheme> = {
	light: {
		background: "#ffffff",
		foreground: "#0a0a0a",
		selectionBackground: "#ebebeb",
		selectionInactiveBackground: "#f5f5f5",
	},
	dark: { background: "#0a0a0a", foreground: "#fbfbfb" },
};

type SearchDecorations = NonNullable<ISearchOptions["decorations"]>;

const XTERM_SEARCH_DECORATIONS: Record<TerminalTheme, SearchDecorations> = {
	light: {
		matchBackground: "#ffe082",
		activeMatchBackground: "#f0c000",
		matchOverviewRuler: "#f0c000",
		activeMatchColorOverviewRuler: "#c09000",
	},
	dark: {
		matchBackground: "#5d4037",
		activeMatchBackground: "#f0c000",
		matchOverviewRuler: "#f0c000",
		activeMatchColorOverviewRuler: "#ffd54f",
	},
};

const openTerminalLink = (_: unknown, uri: string) => {
	externalBrowserService.browserOpenURL(uri);
};

export const xtermTerminal: TerminalFactory = (_commandId, initialTheme) => {
	let theme = initialTheme;
	// Set while a terminal is attached: match decorations are painted once, so a
	// theme change has to re-run the search that painted them.
	let repaintSearch: (() => void) | null = null;

	const terminal = new Terminal({
		allowProposedApi: true,
		convertEol: true,
		scrollback: TERMINAL_SCROLLBACK,
		disableStdin: true,
		fontFamily: "monospace",
		theme: XTERM_THEMES[theme],
	});

	const outputTerminal: OutputTerminal = {
		writeln: (line) => terminal.writeln(line),
		reset: () => terminal.reset(),
		setTheme: (nextTheme) => {
			theme = nextTheme;
			terminal.options.theme = XTERM_THEMES[nextTheme];
			repaintSearch?.();
		},
		dispose: () => terminal.dispose(),

		attach: (element, copyText) => {
			if (terminal.element) {
				// Opened before — re-attaching the existing node keeps the scrollback
				element.appendChild(terminal.element);
			} else {
				terminal.open(element);
			}

			const fit = new FitAddon();
			const links = new WebLinksAddon(openTerminalLink);
			const search = new SearchAddon();

			terminal.loadAddon(fit);
			terminal.loadAddon(links);
			terminal.loadAddon(search);

			fit.fit();

			terminal.attachCustomKeyEventHandler((event) => {
				const isCopyShortcut =
					event.ctrlKey && !event.altKey && event.key.toLowerCase() === "c";
				if (!isCopyShortcut) {
					return true;
				}

				event.preventDefault();
				if (event.type === "keydown" && !event.repeat) {
					const selection = terminal.getSelection();
					if (selection) {
						copyText(selection);
					}
				}
				return false;
			});

			let lastQuery = "";
			repaintSearch = () => {
				if (lastQuery) {
					search.findNext(lastQuery, {
						decorations: XTERM_SEARCH_DECORATIONS[theme],
						incremental: true,
					});
				}
			};

			return {
				fit: () => fit.fit(),

				search: {
					findNext: (query, incremental) => {
						lastQuery = query;
						search.findNext(query, {
							decorations: XTERM_SEARCH_DECORATIONS[theme],
							incremental,
						});
					},
					findPrevious: (query) => {
						lastQuery = query;
						search.findPrevious(query, {
							decorations: XTERM_SEARCH_DECORATIONS[theme],
						});
					},
					clear: () => {
						lastQuery = "";
						search.clearDecorations();
						search.clearActiveDecoration();
					},
					onResults: (listener) => {
						const handle = search.onDidChangeResults(({ resultCount }) =>
							listener(resultCount),
						);
						return () => handle.dispose();
					},
				},
				hasSelection: () => terminal.hasSelection(),
				copySelection: () => {
					const selection = terminal.getSelection();
					if (selection) {
						copyText(selection);
					}
				},

				detach: () => {
					repaintSearch = null;
					terminal.attachCustomKeyEventHandler(() => true);
					search.dispose();
					// Detach the DOM only — the emulator outlives the view
					if (terminal.element && element.contains(terminal.element)) {
						element.removeChild(terminal.element);
					}
					fit.dispose();
					links.dispose();
				},
			};
		},
	};

	return outputTerminal;
};
