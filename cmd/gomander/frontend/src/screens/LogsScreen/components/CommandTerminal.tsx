import { useEffect, useRef, useState } from "react";

import { attachCommandOutput } from "@/commandOutput/commandOutput.ts";
import type { TerminalSearch } from "@/commandOutput/ports.ts";
import { useShortcut } from "@/hooks/useShortcut.ts";
import { useTerminalSearch } from "../hooks/useTerminalSearch.ts";
import { TerminalSearchBar } from "./TerminalSearchBar.tsx";

type Props = {
	commandId: string;
};

export const CommandTerminal = ({ commandId }: Props) => {
	const containerRef = useRef<HTMLDivElement>(null);

	const searchRef = useRef<TerminalSearch | null>(null);

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
	} = useTerminalSearch(searchRef);

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
		if (!containerRef.current) {
			return;
		}
		const container = containerRef.current;

		const terminal = attachCommandOutput(commandId, container);
		searchRef.current = terminal.search;

		const stopListeningResults = terminal.search.onResults(setResultCount);

		const ro = new ResizeObserver(() => terminal.fit());
		ro.observe(container);

		return () => {
			ro.disconnect();
			stopListeningResults();
			searchRef.current = null;
			terminal.detach();
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
