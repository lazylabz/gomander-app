export type TerminalTheme = "light" | "dark";

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
	detach: () => void;
};

export type OutputTerminal = {
	writeln: (line: string) => void;
	reset: () => void;
	setTheme: (theme: TerminalTheme) => void;
	attach: (element: HTMLElement) => AttachedTerminal;
	dispose: () => void;
};

export type TerminalFactory = (
	commandId: string,
	theme: TerminalTheme,
) => OutputTerminal;
