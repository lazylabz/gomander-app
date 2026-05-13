import "@xterm/xterm/css/xterm.css";
import { FitAddon } from "@xterm/addon-fit";
import { SearchAddon } from "@xterm/addon-search";
import { WebLinksAddon } from "@xterm/addon-web-links";
import type { ITheme } from "@xterm/xterm";
import { useEffect, useRef, useState } from "react";

import { useTheme } from "@/contexts/theme.tsx";
import { externalBrowserService } from "@/contracts/service.ts";
import { useShortcut } from "@/hooks/useShortcut.ts";
import { terminalStore } from "@/store/terminalStore.ts";
import { useTerminalSearch } from "../hooks/useTerminalSearch.ts";
import { TerminalSearchBar } from "./TerminalSearchBar.tsx";

export const XTERM_THEMES: Record<"light" | "dark", ITheme> = {
	light: {
		background: "#ffffff",
		foreground: "#0a0a0a",
		selectionBackground: "#ebebeb",
		selectionInactiveBackground: "#f5f5f5",
	},
	dark: { background: "#0a0a0a", foreground: "#fbfbfb" },
};

type Props = {
	commandId: string;
};

const openTerminalLink = (_: unknown, uri: string) => {
	externalBrowserService.browserOpenURL(uri);
};

export const CommandTerminal = ({ commandId }: Props) => {
	const { theme } = useTheme();

	const containerRef = useRef<HTMLDivElement>(null);

	const fitRef = useRef<FitAddon | null>(null);
	const searchRef = useRef<SearchAddon | null>(null);

	const searchInputRef = useRef<HTMLInputElement | null>(null);

	const [searchOpen, setSearchOpen] = useState(false);

	const {
		searchQuery,
		resultCount,
		setResultCount,
		handleSearchInputChange,
		clearSearch,
		nextMatch,
		prevMatch,
	} = useTerminalSearch(searchRef, theme);

	useShortcut("Mod-F", () => setSearchOpen(true));

	useEffect(() => {
		if (searchOpen) {
			searchInputRef.current?.focus();
		}
	}, [searchOpen]);

	const handleClose = () => {
		setSearchOpen(false);
		clearSearch();
	};

	useEffect(() => {
		if (!containerRef.current) return;
		const container = containerRef.current;

		const { getOrCreate, currentTheme } = terminalStore.getState();
		const term = getOrCreate(commandId, currentTheme);

		// Correct theme in case terminal was pre-created by EventListenersContainer
		term.options.theme = currentTheme;

		if (term.element) {
			// Terminal was previously opened — re-attach its DOM element
			container.appendChild(term.element);
		} else {
			term.open(container);
		}

		const fit = new FitAddon();
		const links = new WebLinksAddon(openTerminalLink);
		const search = new SearchAddon();

		term.loadAddon(fit);
		term.loadAddon(links);
		term.loadAddon(search);

		fit.fit();
		fitRef.current = fit;
		searchRef.current = search;

		const resultsHandle = search.onDidChangeResults(
			({ resultCount: count }) => {
				setResultCount(count);
			},
		);

		const ro = new ResizeObserver(() => fitRef.current?.fit());
		ro.observe(container);

		return () => {
			// Disconnect resize observer
			ro.disconnect();

			// Dispose/clean search related vars
			resultsHandle.dispose();
			search.dispose();
			searchRef.current = null;

			// Detach terminal DOM — do NOT dispose, terminal lives in terminalStore
			if (term.element && container.contains(term.element)) {
				container.removeChild(term.element);
			}

			// Dispose/clean fit related vars
			fit.dispose();
			fitRef.current = null;

			// Dispose/clean links addon
			links.dispose();
		};
	}, [commandId, setResultCount]);

	return (
		<div className="absolute inset-0">
			<div ref={containerRef} className="absolute inset-0" />
			{searchOpen && (
				<TerminalSearchBar
					inputRef={searchInputRef}
					query={searchQuery}
					resultCount={resultCount}
					onChange={handleSearchInputChange}
					onPrev={prevMatch}
					onNext={nextMatch}
					onClose={handleClose}
				/>
			)}
		</div>
	);
};
