import type { ITheme } from "@xterm/xterm";
import { Terminal } from "@xterm/xterm";
import { createStore } from "zustand/vanilla";

type TerminalState = {
	terminals: Map<string, Terminal>;
	pendingLogs: Map<string, string[]>;
	currentTheme: ITheme;
	getOrCreate: (commandId: string) => Terminal;
	bufferLogs: (commandId: string, lines: string[]) => void;
	drainPendingLogs: (commandId: string) => string[];
	dispose: (commandId: string) => void;
	disposeAll: () => void;
	setThemeAll: (theme: ITheme) => void;
	resetTerminal: (commandId: string) => void;
};

export const terminalStore = createStore<TerminalState>()((set, get) => ({
	terminals: new Map(),
	pendingLogs: new Map(),
	currentTheme: { background: "#0a0a0a", foreground: "#fbfbfb" }, // dark default

	getOrCreate: (commandId) => {
		const { terminals, currentTheme } = get();
		if (terminals.has(commandId)) {
			return terminals.get(commandId) as Terminal;
		}
		const term = new Terminal({
			allowProposedApi: true,
			convertEol: true,
			scrollback: 10_000,
			disableStdin: true,
			fontFamily: "monospace",
			theme: currentTheme,
		});
		terminals.set(commandId, term);
		set({ terminals: new Map(terminals) });
		return term;
	},

	bufferLogs: (commandId, lines) => {
		const { pendingLogs } = get();
		const existing = pendingLogs.get(commandId) ?? [];
		pendingLogs.set(commandId, [...existing, ...lines]);
		set({ pendingLogs: new Map(pendingLogs) });
	},

	drainPendingLogs: (commandId) => {
		const { pendingLogs } = get();
		const lines = pendingLogs.get(commandId) ?? [];
		pendingLogs.delete(commandId);
		set({ pendingLogs: new Map(pendingLogs) });
		return lines;
	},

	dispose: (commandId) => {
		const { terminals, pendingLogs } = get();
		terminals.get(commandId)?.dispose();
		terminals.delete(commandId);
		pendingLogs.delete(commandId);
		set({ terminals: new Map(terminals), pendingLogs: new Map(pendingLogs) });
	},

	disposeAll: () => {
		const { terminals } = get();
		for (const term of terminals.values()) term.dispose();
		set({ terminals: new Map(), pendingLogs: new Map() });
	},

	setThemeAll: (theme) => {
		const { terminals } = get();
		for (const term of terminals.values()) {
			term.options.theme = theme;
		}
		set({ currentTheme: theme });
	},

	resetTerminal: (commandId: string) => {
		const { terminals } = get();
		terminals.get(commandId)?.reset();
	},
}));
