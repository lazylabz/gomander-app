import { Play, Square } from "lucide-react";

import { useSidebarContext } from "@/components/layout/AppSidebar/contexts/sidebarContext.tsx";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "@/components/ui/context-menu.tsx";
import { SidebarMenuButton } from "@/components/ui/sidebar.tsx";
import { CommandStatus, useDataContext } from "@/contexts/DataContext.tsx";
import { cn } from "@/lib/utils.ts";
import type { Command } from "@/types/contracts.ts";

export const CommandMenuItem = ({ command }: { command: Command }) => {
  const {
    execCommand,
    setActiveCommandId,
    activeCommandId,
    deleteCommand,
    duplicateCommand,
    setCommandStatus,
    commandsStatus,
    stopRunningCommand,
  } = useDataContext();

  const { editCommand } = useSidebarContext();

  const handleRunCommand = async () => {
    setActiveCommandId(command.id);
    setCommandStatus(command.id, CommandStatus.RUNNING);
    await execCommand(command.id);
  };

  const handleDeleteCommand = async () => {
    await deleteCommand(command.id);
    setActiveCommandId(null); // Reset active command after deletion
  };

  const handleEditCommand = () => {
    editCommand(command);
  };

  const handleDuplicateCommand = async () => {
    await duplicateCommand(command);
    setActiveCommandId(null); // Reset active command after duplication
  };

  const onCommandSectionClick = () => {
    setActiveCommandId(command.id);
  };

  const handleStopCommand = async () => {
    await stopRunningCommand(command.id);
  };

  const isIdle = commandsStatus[command.id] === CommandStatus.IDLE;
  const isRunning = commandsStatus[command.id] === CommandStatus.RUNNING;
  const isActiveCommand = activeCommandId === command.id;

  return (
    <ContextMenu>
      <ContextMenuTrigger>
        <SidebarMenuButton
          asChild
          className={cn(
            isActiveCommand && "bg-sidebar-accent",
            isRunning &&
              "bg-green-100 hover:bg-green-100 focus:bg-green-100 active:bg-green-100",
            isActiveCommand &&
              isRunning &&
              "bg-green-200 hover:bg-green-200 focus:bg-green-200 active:bg-green-200",
          )}
        >
          <div
            onClick={onCommandSectionClick}
            className="flex flex-row justify-between items-center w-full"
          >
            {command.name}
            {isIdle && (
              <Play
                size={18}
                className="text-muted-foreground cursor-pointer hover:text-primary"
                onClick={handleRunCommand}
              />
            )}
            {isRunning && (
              <div className="group/command p-0 m-0">
                <Square
                  size={18}
                  className="text-muted-foreground cursor-pointer hover:text-primary"
                  onClick={handleStopCommand}
                />
              </div>
            )}
          </div>
        </SidebarMenuButton>
      </ContextMenuTrigger>
      <ContextMenuContent>
        <ContextMenuItem onClick={handleEditCommand}>Edit</ContextMenuItem>
        <ContextMenuItem onClick={handleDeleteCommand}>Delete</ContextMenuItem>
        <ContextMenuItem onClick={handleDuplicateCommand}>
          Duplicate
        </ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  );
};
