// Holds what was done to each terminal: drive the pipeline and then query
// `terminals.get(commandId)`.
import type {
	OutputTerminal,
	TerminalFactory,
	TerminalTheme,
} from "@/commandOutput/ports.ts";
import { copyTextToClipboard } from "@/useCases/clipboard/copyTextToClipboard.ts";

export type RecordedSearch = {
	query: string;
	direction: "next" | "previous";
	incremental: boolean;
};

export type RecordedTerminal = {
	written: string[];
	resets: number;
	theme: TerminalTheme;
	attachedTo: HTMLElement | null;
	// How many lines had been written when the terminal was first attached.
	writtenWhenAttached: number | null;
	fits: number;
	disposed: boolean;
	searches: RecordedSearch[];
	searchesCleared: number;
	selection: string;
	reportResults: (resultCount: number) => void;
};

export type RecordingTerminals = {
	create: TerminalFactory;
	terminals: Map<string, RecordedTerminal>;
};

export const createRecordingTerminals = (): RecordingTerminals => {
	const terminals = new Map<string, RecordedTerminal>();

	const create: TerminalFactory = (commandId, theme) => {
		const resultListeners: ((resultCount: number) => void)[] = [];

		const recorded: RecordedTerminal = {
			written: [],
			resets: 0,
			theme,
			attachedTo: null,
			writtenWhenAttached: null,
			fits: 0,
			disposed: false,
			searches: [],
			searchesCleared: 0,
			selection: "",
			reportResults: (resultCount) => {
				for (const listener of resultListeners) {
					listener(resultCount);
				}
			},
		};
		terminals.set(commandId, recorded);

		const terminal: OutputTerminal = {
			writeln: (line) => {
				recorded.written.push(line);
			},
			reset: () => {
				recorded.resets += 1;
				recorded.written = [];
			},
			setTheme: (nextTheme) => {
				recorded.theme = nextTheme;
			},
			dispose: () => {
				recorded.disposed = true;
			},
			attach: (element) => {
				recorded.attachedTo = element;
				recorded.writtenWhenAttached ??= recorded.written.length;
				return {
					fit: () => {
						recorded.fits += 1;
					},
					search: {
						findNext: (query, incremental = false) => {
							recorded.searches.push({
								query,
								direction: "next",
								incremental,
							});
						},
						findPrevious: (query) => {
							recorded.searches.push({
								query,
								direction: "previous",
								incremental: false,
							});
						},
						clear: () => {
							recorded.searchesCleared += 1;
						},
						onResults: (listener) => {
							resultListeners.push(listener);
							return () => {
								resultListeners.splice(resultListeners.indexOf(listener), 1);
							};
						},
					},
					hasSelection: () => recorded.selection.length > 0,
					copySelection: () => {
						if (recorded.selection) {
							void copyTextToClipboard(recorded.selection);
						}
					},
					detach: () => {
						recorded.attachedTo = null;
					},
				};
			},
		};

		return terminal;
	};

	return { create, terminals };
};
