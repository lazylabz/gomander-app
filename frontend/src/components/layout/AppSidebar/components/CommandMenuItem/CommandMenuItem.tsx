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
import type { Command } from "@/types/contracts.ts";

export const CommandMenuItem = ({ command }: { command: Command }) => {
  const {
    execCommand,
    setActiveCommandId,
    activeCommandId,
    deleteCommand,
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

  const onCommandSectionClick = () => {
    setActiveCommandId(command.id);
  };

  const handleStopCommand = async () => {
    await stopRunningCommand(command.id);
  };

  return (
    <ContextMenu>
      <ContextMenuTrigger>
        <SidebarMenuButton asChild isActive={activeCommandId === command.id}>
          <div
            onClick={onCommandSectionClick}
            className="flex flex-row justify-between items-center w-full"
          >
            {command.name}
            {commandsStatus[command.id] === CommandStatus.IDLE && (
              <Play
                size={18}
                className="text-muted-foreground cursor-pointer hover:text-primary"
                onClick={handleRunCommand}
              />
            )}
            {commandsStatus[command.id] === CommandStatus.RUNNING && (
              <div className="group/command p-0 m-0">
                <Play
                  size={18}
                  fill="currentColor"
                  className="text-green-600 group-hover/command:hidden"
                  onClick={handleRunCommand}
                />
                <Square
                  size={18}
                  fill="currentColor"
                  className="text-red-400 hidden group-hover/command:block cursor-pointer"
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
      </ContextMenuContent>
    </ContextMenu>
  );
};
