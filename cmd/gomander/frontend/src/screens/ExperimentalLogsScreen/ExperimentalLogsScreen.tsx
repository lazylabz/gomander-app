import {
	SidebarTrigger,
	useSidebar,
} from "@/design-system/components/ui/sidebar.tsx";
import { cn } from "@/design-system/lib/utils.ts";
import { useCommandStore } from "@/store/commandStore.ts";
import { CommandTerminal } from "./CommandTerminal.tsx";

export const ExperimentalLogsScreen = () => {
	const activeCommandId = useCommandStore((s) => s.activeCommandId);
	const { open, isMobile } = useSidebar();
	const isDesktopSidebarOpen = !isMobile && open;

	return (
		<div className="relative h-full w-full bg-background overflow-hidden flex flex-col p-2">
			<div
				className={cn(
					"fixed top-3.5 z-10 bg-background rounded-sm duration-200 ease-linear opacity-50 hover:opacity-100 transition-all",
					isDesktopSidebarOpen ? "left-[16.5rem]" : "left-2",
				)}
			>
				<SidebarTrigger />
			</div>
			<div className="relative flex-1 min-h-0 w-full">
				{activeCommandId && (
					<CommandTerminal key={activeCommandId} commandId={activeCommandId} />
				)}
			</div>
		</div>
	);
};
