import "@xterm/xterm/css/xterm.css";
import { FitAddon } from "@xterm/addon-fit";
import type { ITheme } from "@xterm/xterm";
import { useEffect, useRef } from "react";

import { useTheme } from "@/contexts/theme.tsx";
import { commandStore } from "@/store/commandStore.ts";
import { terminalStore } from "@/store/terminalStore.ts";

export const XTERM_THEMES: Record<"light" | "dark", ITheme> = {
	light: { background: "#ffffff", foreground: "#0a0a0a" },
	dark: { background: "#0a0a0a", foreground: "#fbfbfb" },
};

type Props = {
	commandId: string;
};

export const CommandTerminal = ({ commandId }: Props) => {
	const { theme } = useTheme();
	const containerRef = useRef<HTMLDivElement>(null);
	const fitRef = useRef<FitAddon | null>(null);
	const xtermThemeRef = useRef(XTERM_THEMES[theme]);
	xtermThemeRef.current = XTERM_THEMES[theme];

	// Mount / re-attach terminal on activation; detach (without dispose) on deactivation
	useEffect(() => {
		if (!containerRef.current) return;
		const container = containerRef.current;

		const term = terminalStore
			.getState()
			.getOrCreate(commandId, xtermThemeRef.current);

		if (term.element) {
			// Terminal was previously opened — re-attach its DOM element
			container.appendChild(term.element);
		} else {
			// First open: replay logs currently in store, then open into the container
			const existingLogs =
				commandStore.getState().commandsLogs[commandId] ?? [];
			for (const line of existingLogs) term.writeln(line);
			term.open(container);
		}

		const fit = new FitAddon();
		term.loadAddon(fit);
		fit.fit();
		fitRef.current = fit;

		const ro = new ResizeObserver(() => fitRef.current?.fit());
		ro.observe(container);

		return () => {
			ro.disconnect();
			// Detach terminal DOM — do NOT dispose, terminal lives in terminalStore
			if (term.element && container.contains(term.element)) {
				container.removeChild(term.element);
			}
			fit.dispose();
			fitRef.current = null;
		};
	}, [commandId]);

	// Propagate theme change to all open terminals
	useEffect(() => {
		terminalStore.getState().setThemeAll(XTERM_THEMES[theme]);
	}, [theme]);

	return <div ref={containerRef} className="absolute inset-0" />;
};
