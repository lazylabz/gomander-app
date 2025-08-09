import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { GripVertical, Play, Square } from "lucide-react";

import { useSidebarContext } from "@/components/layout/AppSidebarLayout/components/AppSidebar/contexts/sidebarContext.tsx";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "@/components/ui/context-menu.tsx";
import { SidebarMenuButton } from "@/components/ui/sidebar.tsx";
import type { Command } from "@/contracts/types.ts";
import { cn } from "@/lib/utils.ts";
import { useCommandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { deleteCommand } from "@/useCases/command/deleteCommand.ts";
import { duplicateCommand } from "@/useCases/command/duplicateCommand.ts";
import { startCommand } from "@/useCases/command/startCommand.ts";
import { stopCommand } from "@/useCases/command/stopCommand.ts";

export const CommandMenuItem = ({
  command,
  draggable = false,
}: {
  command: Command;
  draggable?: boolean;
}) => {
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
    await startCommand(command.id);
  };

  const handleDeleteCommand = async () => {
    await deleteCommand(command.id);
    setActiveCommandId(null); // Reset active command after deletion
  };

  const handleEditCommand = () => {
    startEditingCommand(command);
  };

  const handleDuplicateCommand = async () => {
    await duplicateCommand(command);
    setActiveCommandId(null); // Reset active command after duplication
  };

  const onCommandSectionClick = () => {
    setActiveCommandId(command.id);
  };

  const handleStopCommand = async () => {
    await stopCommand(command.id);
    setActiveCommandId(command.id);
  };

  const isIdle = commandsStatus[command.id] === CommandStatus.IDLE;
  const isRunning = commandsStatus[command.id] === CommandStatus.RUNNING;
  const isActiveCommand = activeCommandId === command.id;

  const className = cn(
    isActiveCommand && "bg-sidebar-accent",
    isRunning &&
      "bg-green-100 hover:bg-green-100 focus:bg-green-100 active:bg-green-100",
    isActiveCommand &&
      isRunning &&
      "bg-green-200 hover:bg-green-200 focus:bg-green-200 active:bg-green-200",
  );

  return (
    <ContextMenu>
      <ContextMenuTrigger>
        <SidebarMenuButton asChild className={className}>
          <div
            onClick={onCommandSectionClick}
            className="flex flex-row justify-between items-center w-full"
            ref={setNodeRef}
            style={style}
          >
            <div className="flex items-center gap-1 w-full text-sm text-sidebar-foreground">
              {draggable && (
                <div
                  {...attributes}
                  {...listeners}
                  className="cursor-grab active:cursor-grabbing pr-0.5 rounded hover:bg-sidebar-accent/50"
                >
                  <GripVertical size={14} className="text-muted-foreground" />
                </div>
              )}
              {command.name}
            </div>
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
        <ContextMenuItem disabled={isRunning} onClick={handleEditCommand}>
          Edit
        </ContextMenuItem>
        <ContextMenuItem disabled={isRunning} onClick={handleDeleteCommand}>
          Delete
        </ContextMenuItem>
        <ContextMenuItem disabled={isRunning} onClick={handleDuplicateCommand}>
          Duplicate
        </ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  );
};
