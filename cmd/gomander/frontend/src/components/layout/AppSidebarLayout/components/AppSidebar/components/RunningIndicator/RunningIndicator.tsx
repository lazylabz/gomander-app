import { useTranslation } from "react-i18next";

import {
	Tooltip,
	TooltipContent,
	TooltipTrigger,
} from "@/design-system/components/ui/tooltip.tsx";
import { cn } from "@/design-system/lib/utils.ts";
import { useCommandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

export const RunningIndicator = ({ className }: { className?: string }) => {
	const { t } = useTranslation();
	const commands = useCommandStore((state) => state.commands);
	const commandsStatus = useCommandStore((state) => state.commandsStatus);

	const runningCount = commands.filter(
		(command) => commandsStatus[command.id] === CommandStatus.RUNNING,
	).length;

	const isRunning = runningCount > 0;

	const label = isRunning
		? t("sidebar.runningIndicator.running", { count: runningCount })
		: t("sidebar.runningIndicator.idle");

	return (
		<Tooltip>
			<TooltipTrigger
				className={cn("flex items-center justify-center", className)}
			>
				<span className="relative flex size-1.5">
					{isRunning && (
						<span className="absolute inline-flex size-full animate-ping rounded-full bg-green-500 opacity-75" />
					)}
					<span
						className={cn(
							"relative inline-flex size-full rounded-full transition-colors",
							isRunning ? "bg-green-500" : "bg-muted-foreground/40",
						)}
					/>
				</span>
			</TooltipTrigger>
			<TooltipContent>{label}</TooltipContent>
		</Tooltip>
	);
};
