import { type ChangeEvent, type RefObject, useEffect, useState } from "react";

import type { TerminalSearch } from "@/commandOutput/ports.ts";

export const useTerminalSearch = (
	searchRef: RefObject<TerminalSearch | null>,
) => {
	const [searchQuery, setSearchQuery] = useState("");
	const [resultCount, setResultCount] = useState(0);
	useEffect(() => {
		if (!searchRef.current) {
			return;
		}
		if (!searchQuery) {
			searchRef.current.clear();
			setResultCount(0);
			return;
		}
		searchRef.current.findNext(searchQuery, true);
	}, [searchQuery, searchRef]);

	const clearSearch = () => {
		setSearchQuery("");
		setResultCount(0);
		searchRef.current?.clear();
	};

	const nextMatch = () => {
		if (!searchRef.current || !searchQuery) {
			return;
		}
		searchRef.current.findNext(searchQuery);
	};

	const prevMatch = () => {
		if (!searchRef.current || !searchQuery) {
			return;
		}
		searchRef.current.findPrevious(searchQuery);
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
