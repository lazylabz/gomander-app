import { ChevronLeft, ChevronRight, X } from "lucide-react";
import type { ChangeEvent, KeyboardEvent, RefObject } from "react";
import { useTranslation } from "react-i18next";

import { Input } from "@/design-system/components/ui/input.tsx";

type Props = {
	inputRef: RefObject<HTMLInputElement | null>;
	query: string;
	resultCount: number;
	onChange: (e: ChangeEvent<HTMLInputElement>) => void;
	onPrev: () => void;
	onNext: () => void;
	onClose: () => void;
};

export const TerminalSearchBar = ({
	inputRef,
	query,
	resultCount,
	onChange,
	onPrev,
	onNext,
	onClose,
}: Props) => {
	const { t } = useTranslation();

	const handleKeyDown = (event: KeyboardEvent<HTMLInputElement>) => {
		if (event.key === "Escape") {
			onClose();
		} else if (
			event.key === "ArrowDown" ||
			(event.key === "Enter" && !event.shiftKey)
		) {
			event.preventDefault();
			onNext();
		} else if (
			event.key === "ArrowUp" ||
			(event.key === "Enter" && event.shiftKey)
		) {
			event.preventDefault();
			onPrev();
		}
	};

	return (
		<div className="absolute top-2 right-2 z-10 flex flex-col bg-background gap-1.5">
			<Input
				ref={inputRef}
				autoComplete="off"
				autoCorrect="off"
				autoCapitalize="off"
				className="w-64"
				value={query}
				onChange={onChange}
				onKeyDown={handleKeyDown}
			/>
			<div className="text-xs text-muted-foreground pl-2 flex items-center gap-2 pb-1 justify-between select-none">
				<div className="flex flex-row items-center gap-2">
					<div className="flex flex-row items-center">
						<button
							type="button"
							onClick={onPrev}
							aria-label={t("logs.searchPrev")}
							className="inline-flex items-center border-0 bg-transparent p-0 text-muted-foreground hover:text-foreground cursor-pointer"
						>
							<ChevronLeft size={14} />
						</button>
						<button
							type="button"
							onClick={onNext}
							aria-label={t("logs.searchNext")}
							className="inline-flex items-center border-0 bg-transparent p-0 text-muted-foreground hover:text-foreground cursor-pointer"
						>
							<ChevronRight size={14} />
						</button>
					</div>
					{t("logs.matches", { count: resultCount })}
				</div>
				<button
					type="button"
					onClick={onClose}
					aria-label={t("common.close")}
					className="inline-flex items-center border-0 bg-transparent p-0 text-muted-foreground hover:text-foreground cursor-pointer"
				>
					<X size={14} />
				</button>
			</div>
		</div>
	);
};
