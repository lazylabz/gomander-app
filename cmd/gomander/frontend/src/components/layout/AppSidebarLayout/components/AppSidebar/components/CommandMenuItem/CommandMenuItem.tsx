import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { GripVertical, LinkIcon, Play, Square } from "lucide-react";
import { useTranslation } from "react-i18next";

import { useSidebarContext } from "@/components/layout/AppSidebarLayout/components/AppSidebar/contexts/sidebarContext.tsx";
import { useTheme } from "@/contexts/theme.tsx";
import { externalBrowserService } from "@/contracts/service.ts";
import type { Command } from "@/contracts/types.ts";
import {
	ContextMenu,
	ContextMenuContent,
	ContextMenuItem,
	ContextMenuTrigger,
} from "@/design-system/components/ui/context-menu.tsx";
import { SidebarMenuButton } from "@/design-system/components/ui/sidebar.tsx";
import { cn } from "@/design-system/lib/utils.ts";
import { useCommandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { cleanCommandError } from "@/useCases/command/cleanCommandError.ts";
import { deleteCommand } from "@/useCases/command/deleteCommand.ts";
import { duplicateCommand } from "@/useCases/command/duplicateCommand.ts";
import { removeCommandFromGroup } from "@/useCases/command/removeCommandFromGroup.ts";
import { startCommand } from "@/useCases/command/startCommand.ts";
import { stopCommand } from "@/useCases/command/stopCommand.ts";

// Each control sits above the row's button, which stretches across the row, so
// a control is its own focusable button rather than a child the row swallows.
const controlClassName =
	"relative z-10 flex items-center justify-center rounded-sm outline-hidden ring-sidebar-ring focus-visible:ring-2";

export const CommandMenuItem = ({
	command,
	insideGroupId,
	draggable = false,
}: {
	command: Command;
	insideGroupId?: string;
	draggable?: boolean;
}) => {
	const { t } = useTranslation();
	const { theme } = useTheme();

	const setActiveCommandId = useCommandStore(
		(state) => state.setActiveCommandId,
	);
	const commandsStatus = useCommandStore((state) => state.commandsStatus);
	const activeCommandId = useCommandStore((state) => state.activeCommandId);
	const commandIdsWithErrors = useCommandStore(
		(state) => state.commandIdsWithErrors,
	);

	const { startEditingCommand } = useSidebarContext();

	const { attributes, listeners, setNodeRef, transform, transition } =
		useSortable({ id: command.id });

	const style = {
		transform: CSS.Transform.toString(transform),
		transition,
	};

	const selectCommand = () => {
		setActiveCommandId(command.id);
		cleanCommandError(command.id);
	};

	const handleRunCommand = async () => {
		selectCommand();
		await startCommand(command.id);
	};

	const handleOpenLink = () => {
		if (!command.link) {
			return;
		}

		selectCommand();
		externalBrowserService.browserOpenURL(command.link);
	};

	const handleDeleteCommand = async () => {
		await deleteCommand(command.id);
	};

	const handleRemoveFromGroup = async () => {
		if (!insideGroupId) {
			return;
		}

		await removeCommandFromGroup(command.id, insideGroupId);
	};

	const handleEditCommand = () => {
		startEditingCommand(command);
	};

	const handleDuplicateCommand = async () => {
		await duplicateCommand(command.id, insideGroupId);
	};

	const handleStopCommand = async () => {
		await stopCommand(command.id);
		selectCommand();
	};

	const isIdle = commandsStatus[command.id] === CommandStatus.IDLE;
	const isRunning = commandsStatus[command.id] === CommandStatus.RUNNING;
	const isActiveCommand = activeCommandId === command.id;
	const hasError = commandIdsWithErrors.includes(command.id);

	const className = cn(
		"p-2.5",
		isActiveCommand && "bg-sidebar-accent",
		isRunning &&
			"bg-green-100 hover:bg-green-200 focus:bg-green-100 active:bg-green-100",
		isRunning &&
			theme === "dark" &&
			"bg-green-300/30 hover:bg-green-200/40 focus:bg-green-300/30 active:bg-green-300/30",
		isActiveCommand &&
			isRunning &&
			"bg-green-200 hover:bg-green-200 focus:bg-green-200 active:bg-green-200",
		isActiveCommand &&
			isRunning &&
			theme === "dark" &&
			"bg-green-200/40 hover:bg-green-200/40 focus:bg-green-200/40 active:bg-green-200/40",
		hasError &&
			"bg-red-100 hover:bg-red-200 focus:bg-red-100 active:bg-red-100",
		hasError &&
			theme === "dark" &&
			"bg-red-300/30 hover:bg-red-200/40 focus:bg-red-300/30 active:bg-red-300/30",
	);

	return (
		<ContextMenu>
			<ContextMenuTrigger>
				<SidebarMenuButton asChild className={className}>
					<div
						className="relative flex flex-row justify-between items-center w-full select-none"
						ref={setNodeRef}
						style={style}
					>
						<div
							className={cn(
								"flex items-center gap-1 w-full min-w-0 overflow-hidden text-sm text-sidebar-foreground",
								!draggable && "pl-2",
							)}
						>
							{draggable && (
								<div
									{...attributes}
									{...listeners}
									className={cn(
										controlClassName,
										"cursor-grab active:cursor-grabbing pr-0.5 rounded hover:bg-sidebar-accent/50 shrink-0",
									)}
								>
									<GripVertical size={14} className="text-muted-foreground" />
								</div>
							)}
							<button
								type="button"
								onClick={selectCommand}
								title={command.name}
								className="cursor-default text-left truncate outline-hidden after:absolute after:inset-0 after:rounded-md after:ring-sidebar-ring focus-visible:after:ring-2"
							>
								{command.name}
							</button>
							{command.link && (
								<button
									type="button"
									aria-label={t("sidebar.commands.openLink")}
									onClick={handleOpenLink}
									className={cn(
										controlClassName,
										"text-muted-foreground cursor-pointer hover:text-primary min-w-4",
									)}
								>
									<LinkIcon size={12} />
								</button>
							)}
						</div>
						<div className="shrink-0 flex">
							{isIdle && (
								<button
									type="button"
									aria-label={t("sidebar.commands.run")}
									onClick={handleRunCommand}
									className={cn(
										controlClassName,
										"text-muted-foreground cursor-pointer hover:text-primary",
									)}
								>
									<Play size={16} />
								</button>
							)}
							{isRunning && (
								<button
									type="button"
									aria-label={t("sidebar.commands.stop")}
									onClick={handleStopCommand}
									className={cn(
										controlClassName,
										"text-muted-foreground dark:text-primary/70 cursor-pointer hover:text-primary dark:hover:text-primary",
									)}
								>
									<Square size={16} />
								</button>
							)}
						</div>
					</div>
				</SidebarMenuButton>
			</ContextMenuTrigger>
			<ContextMenuContent>
				<ContextMenuItem disabled={isRunning} onClick={handleEditCommand}>
					{t("common.edit")}
				</ContextMenuItem>
				<ContextMenuItem disabled={isRunning} onClick={handleDuplicateCommand}>
					{t("common.duplicate")}
				</ContextMenuItem>
				{!insideGroupId && (
					<ContextMenuItem disabled={isRunning} onClick={handleDeleteCommand}>
						{t("common.delete")}
					</ContextMenuItem>
				)}
				{insideGroupId && (
					<ContextMenuItem disabled={isRunning} onClick={handleRemoveFromGroup}>
						{t("sidebar.commands.removeFromGroup")}
					</ContextMenuItem>
				)}
			</ContextMenuContent>
		</ContextMenu>
	);
};
