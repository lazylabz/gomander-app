import type { ISearchOptions, SearchAddon } from "@xterm/addon-search";
import { type ChangeEvent, type RefObject, useEffect, useState } from "react";

type SearchDecorations = NonNullable<ISearchOptions["decorations"]>;

export const XTERM_SEARCH_DECORATIONS: Record<
	"light" | "dark",
	SearchDecorations
> = {
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

export const useTerminalSearch = (
	searchRef: RefObject<SearchAddon | null>,
	theme: "light" | "dark",
) => {
	const [searchQuery, setSearchQuery] = useState("");
	const [resultCount, setResultCount] = useState(0);
	useEffect(() => {
		if (!searchRef.current) return;
		if (!searchQuery) {
			searchRef.current.clearDecorations();
			setResultCount(0);
			return;
		}
		searchRef.current.findNext(searchQuery, {
			decorations: XTERM_SEARCH_DECORATIONS[theme],
			incremental: true,
		});
	}, [searchQuery, theme, searchRef]);

	const clearSearch = () => {
		setSearchQuery("");
		setResultCount(0);
		searchRef.current?.clearDecorations();
		searchRef.current?.clearActiveDecoration();
	};

	const nextMatch = () => {
		if (!searchRef.current || !searchQuery) return;
		searchRef.current.findNext(searchQuery, {
			decorations: XTERM_SEARCH_DECORATIONS[theme],
		});
	};

	const prevMatch = () => {
		if (!searchRef.current || !searchQuery) return;
		searchRef.current.findPrevious(searchQuery, {
			decorations: XTERM_SEARCH_DECORATIONS[theme],
		});
	};

	const handleSearchInputChange = (event: ChangeEvent<HTMLInputElement>) => {
		setSearchQuery(event.target.value);
	};

	return {
		searchQuery,
		resultCount,
		setResultCount,
		handleSearchInputChange,
		clearSearch,
		nextMatch,
		prevMatch,
	};
};
