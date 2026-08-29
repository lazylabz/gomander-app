import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";

import { attachCommandOutput } from "@/commandOutput/commandOutput.ts";
import type {
	AttachedTerminal,
	TerminalSearch,
} from "@/commandOutput/ports.ts";
import {
	ContextMenu,
	ContextMenuContent,
	ContextMenuItem,
	ContextMenuTrigger,
} from "@/design-system/components/ui/context-menu.tsx";
import { useShortcut } from "@/hooks/useShortcut.ts";
import { copyTextToClipboard } from "@/useCases/clipboard/copyTextToClipboard.ts";
import { useTerminalSearch } from "../hooks/useTerminalSearch.ts";
import { TerminalSearchBar } from "./TerminalSearchBar.tsx";

type Props = {
	commandId: string;
};

export const CommandTerminal = ({ commandId }: Props) => {
	const { t } = useTranslation();
	const containerRef = useRef<HTMLDivElement>(null);
	const attachedRef = useRef<AttachedTerminal | null>(null);

	const searchRef = useRef<TerminalSearch | null>(null);

	const searchInputRef = useRef<HTMLInputElement | null>(null);

	const [searchOpen, setSearchOpen] = useState(false);
	const [canCopy, setCanCopy] = useState(false);

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

	const handleContextMenuOpenChange = (open: boolean) => {
		if (open) {
			setCanCopy(attachedRef.current?.hasSelection() ?? false);
		}
	};

	const handleCopy = () => attachedRef.current?.copySelection();

	useEffect(() => {
		if (!containerRef.current) {
			return;
		}
		const container = containerRef.current;

		const terminal = attachCommandOutput(commandId, container, (text) => {
			void copyTextToClipboard(text);
		});
		attachedRef.current = terminal;
		searchRef.current = terminal.search;

		const stopListeningResults = terminal.search.onResults(setResultCount);

		const ro = new ResizeObserver(() => terminal.fit());
		ro.observe(container);

		return () => {
			ro.disconnect();
			stopListeningResults();
			attachedRef.current = null;
			searchRef.current = null;
			terminal.detach();
		};
	}, [commandId, setResultCount]);

	return (
		<div className="absolute inset-0">
			<ContextMenu onOpenChange={handleContextMenuOpenChange}>
				<ContextMenuTrigger asChild>
					<div ref={containerRef} className="absolute inset-0" />
				</ContextMenuTrigger>
				<ContextMenuContent>
					<ContextMenuItem disabled={!canCopy} onClick={handleCopy}>
						{t("common.copy")}
					</ContextMenuItem>
				</ContextMenuContent>
			</ContextMenu>
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
