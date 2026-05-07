import type { ITheme } from "@xterm/xterm";
import { Terminal } from "@xterm/xterm";
import { createStore } from "zustand/vanilla";

type TerminalState = {
	terminals: Map<string, Terminal>;
	currentTheme: ITheme;
	getOrCreate: (commandId: string, theme: ITheme) => Terminal;
	dispose: (commandId: string) => void;
	disposeAll: () => void;
	setThemeAll: (theme: ITheme) => void;
};

export const terminalStore = createStore<TerminalState>()((set, get) => ({
	terminals: new Map(),
	currentTheme: { background: "#0a0a0a", foreground: "#fbfbfb" }, // dark default

	getOrCreate: (commandId, theme) => {
		const { terminals } = get();
		if (terminals.has(commandId)) {
			return terminals.get(commandId) as Terminal;
		}
		const term = new Terminal({
			allowProposedApi: true,
			convertEol: true,
			scrollback: 10_000,
			disableStdin: true,
			fontFamily: "monospace",
			theme,
		});
		terminals.set(commandId, term);
		set({ terminals });
		return term;
	},

	dispose: (commandId) => {
		const { terminals } = get();
		terminals.get(commandId)?.dispose();
		terminals.delete(commandId);
		set({ terminals });
	},

	disposeAll: () => {
		const { terminals } = get();
		for (const term of terminals.values()) term.dispose();
		terminals.clear();
		set({ terminals });
	},

	setThemeAll: (theme) => {
		const { terminals } = get();
		for (const term of terminals.values()) {
			term.options.theme = theme;
		}
		set({ currentTheme: theme });
	},
}));
