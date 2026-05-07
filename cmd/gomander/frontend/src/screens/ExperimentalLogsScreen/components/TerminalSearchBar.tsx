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
			<span className="text-xs text-muted-foreground pl-2 flex items-center gap-2 pb-1 justify-between select-none">
				<div className="flex flex-row items-center gap-2">
					<div className="flex flex-row items-center">
						<ChevronLeft
							className="text-muted-foreground hover:text-foreground cursor-pointer"
							onClick={onPrev}
							size={14}
						/>
						<ChevronRight
							className="text-muted-foreground hover:text-foreground cursor-pointer"
							onClick={onNext}
							size={14}
						/>
					</div>
					{t("logs.matches", { count: resultCount })}
				</div>
				<X
					size={14}
					onClick={onClose}
					className="text-muted-foreground hover:text-foreground cursor-pointer"
				/>
			</span>
		</div>
	);
};
