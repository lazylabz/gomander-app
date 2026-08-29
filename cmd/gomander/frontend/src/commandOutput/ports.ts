export type TerminalTheme = "light" | "dark";

// How many lines a terminal keeps. Shared with the pipeline: anything it holds
// for a terminal that does not exist yet beyond this would be dropped on write.
export const TERMINAL_SCROLLBACK = 10_000;

export type CopyText = (text: string) => void;

export type TerminalSearch = {
	findNext: (query: string, incremental?: boolean) => void;
	findPrevious: (query: string) => void;
	clear: () => void;
	onResults: (listener: (resultCount: number) => void) => () => void;
};

// What a mounted terminal hands back to the view: sizing, the search the search
// bar drives, and a teardown that leaves the emulator alive off-screen.
export type AttachedTerminal = {
	fit: () => void;
	search: TerminalSearch;
	hasSelection: () => boolean;
	copySelection: () => void;
	detach: () => void;
};

export type OutputTerminal = {
	writeln: (line: string) => void;
	reset: () => void;
	setTheme: (theme: TerminalTheme) => void;
	attach: (element: HTMLElement, copyText: CopyText) => AttachedTerminal;
	dispose: () => void;
};

export type TerminalFactory = (
	commandId: string,
	theme: TerminalTheme,
) => OutputTerminal;
