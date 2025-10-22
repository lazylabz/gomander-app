import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { GripVertical, LinkIcon, Play, Square } from "lucide-react";
import { toast } from "sonner";

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
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { useCommandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { deleteCommand } from "@/useCases/command/deleteCommand.ts";
import { duplicateCommand } from "@/useCases/command/duplicateCommand.ts";
import { removeCommandFromGroup } from "@/useCases/command/removeCommandFromGroup.ts";
import { startCommand } from "@/useCases/command/startCommand.ts";
import { stopCommand } from "@/useCases/command/stopCommand.ts";

export const CommandMenuItem = ({
  command,
  insideGroupId,
  draggable = false,
}: {
  command: Command;
  insideGroupId?: string;
  draggable?: boolean;
}) => {
  const { theme } = useTheme();

  const setActiveCommandId = useCommandStore(
    (state) => state.setActiveCommandId,
  );
  const commandsStatus = useCommandStore((state) => state.commandsStatus);
  const activeCommandId = useCommandStore((state) => state.activeCommandId);

  const { startEditingCommand } = useSidebarContext();

  const { attributes, listeners, setNodeRef, transform, transition } =
    useSortable({ id: command.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  const handleRunCommand = async () => {
    setActiveCommandId(command.id);
    try {
      await startCommand(command.id);
    } catch (e) {
      toast.error(parseError(e, "Failed to run command"));
    }
  };

  const handleOpenLink = async () => {
    if (!command.link) {
      return;
    }

    externalBrowserService.browserOpenURL(command.link);
  };

  const handleDeleteCommand = async () => {
    try {
      await deleteCommand(command.id);
      setActiveCommandId(null); // Reset active command after deletion
      toast.success("Command deleted successfully");
    } catch (e) {
      toast.error(parseError(e, "Failed to delete command"));
    } finally {
      fetchCommands();
      fetchCommandGroups();
    }
    setActiveCommandId(null); // Reset active command after deletion
  };

  const handleRemoveFromGroup = async () => {
    if (!insideGroupId) return;
    try {
      await removeCommandFromGroup(command.id, insideGroupId);
      toast.success("Command removed from group successfully");
    } catch (e) {
      toast.error(parseError(e, "Failed to remove command from group"));
    } finally {
      fetchCommandGroups();
    }
  };

  const handleEditCommand = () => {
    startEditingCommand(command);
  };

  const handleDuplicateCommand = async () => {
    try {
      await duplicateCommand(command, insideGroupId);
      toast.success("Command duplicated successfully");
    } catch (e) {
      toast.error(parseError(e, "Failed to duplicate command"));
    } finally {
      fetchCommands();
      if (insideGroupId) {
        fetchCommandGroups();
      }
    }
    setActiveCommandId(null); // Reset active command after duplication
  };

  const onCommandSectionClick = () => {
    setActiveCommandId(command.id);
  };

  const handleStopCommand = async () => {
    try {
      await stopCommand(command.id);
    } catch (e) {
      toast.error(parseError(e, "Failed to stop command"));
    }
    setActiveCommandId(command.id);
  };

  const isIdle = commandsStatus[command.id] === CommandStatus.IDLE;
  const isRunning = commandsStatus[command.id] === CommandStatus.RUNNING;
  const isActiveCommand = activeCommandId === command.id;

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
  );

  return (
    <ContextMenu>
      <ContextMenuTrigger>
        <SidebarMenuButton asChild className={className}>
          <div
            onClick={onCommandSectionClick}
            className="flex flex-row justify-between items-center w-full select-none"
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
                  className="cursor-grab active:cursor-grabbing pr-0.5 rounded hover:bg-sidebar-accent/50 shrink-0"
                >
                  <GripVertical size={14} className="text-muted-foreground" />
                </div>
              )}
              <p
                className="cursor-default text-left truncate"
                title={command.name}
              >
                {command.name}
              </p>
              {command.link && (
                <LinkIcon
                  size={12}
                  onClick={handleOpenLink}
                  className="text-muted-foreground cursor-pointer hover:text-primary min-w-4"
                />
              )}
            </div>
            <div className="shrink-0">
              {isIdle && (
                <Play
                  size={16}
                  className="text-muted-foreground cursor-pointer hover:text-primary"
                  onClick={handleRunCommand}
                />
              )}
              {isRunning && (
                <Square
                  size={16}
                  className="text-muted-foreground dark:text-primary/70 cursor-pointer hover:text-primary dark:hover:text-primary"
                  onClick={handleStopCommand}
                />
              )}
            </div>
          </div>
        </SidebarMenuButton>
      </ContextMenuTrigger>
      <ContextMenuContent>
        <ContextMenuItem disabled={isRunning} onClick={handleEditCommand}>
          Edit
        </ContextMenuItem>
        <ContextMenuItem disabled={isRunning} onClick={handleDuplicateCommand}>
          Duplicate
        </ContextMenuItem>
        {!insideGroupId && (
          <ContextMenuItem disabled={isRunning} onClick={handleDeleteCommand}>
            Delete
          </ContextMenuItem>
        )}
        {insideGroupId && (
          <ContextMenuItem disabled={isRunning} onClick={handleRemoveFromGroup}>
            Remove from group
          </ContextMenuItem>
        )}
      </ContextMenuContent>
    </ContextMenu>
  );
};
