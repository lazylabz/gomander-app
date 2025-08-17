import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { Folder, FolderOpen, Play, Square } from "lucide-react";
import { type SyntheticEvent, useState } from "react";

import { CommandMenuItem } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/CommandMenuItem/CommandMenuItem.tsx";
import { useSidebarContext } from "@/components/layout/AppSidebarLayout/components/AppSidebar/contexts/sidebarContext.tsx";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "@/components/ui/context-menu.tsx";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuItem,
} from "@/components/ui/sidebar.tsx";
import type { Command, CommandGroup } from "@/contracts/types.ts";
import { useCommandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { deleteCommandGroup } from "@/useCases/commandGroup/deleteCommandGroup.ts";
import { runCommandGroup } from "@/useCases/commandGroup/runCommandGroup.ts";
import { stopCommandGroup } from "@/useCases/commandGroup/stopCommandGroup.ts";

export const CommandGroupSection = ({
  commandGroup,
}: {
  commandGroup: CommandGroup;
}) => {
  const commandsStatus = useCommandStore((state) => state.commandsStatus);

  const { startEditingCommandGroup, isReorderingGroups } = useSidebarContext();

  const [internalIsOpen, setInternalIsOpen] = useState(false);
  const isOpen = internalIsOpen && !isReorderingGroups;

  const { attributes, listeners, setNodeRef, transform, transition } =
    useSortable({ id: commandGroup.id, disabled: !isReorderingGroups });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  const numberOfCommandsRunning = commandGroup.commands.filter(
    (command: Command) => commandsStatus[command.id] === CommandStatus.RUNNING,
  ).length;

  const someCommandIsRunning = numberOfCommandsRunning > 0;

  const someCommandIsIdle = commandGroup.commands.some(
    (command: Command) => commandsStatus[command.id] === CommandStatus.IDLE,
  );

  const run = (e: SyntheticEvent) => {
    // Prevent the folder from collapsing when clicking the play button
    e.stopPropagation();

    if (isReorderingGroups) return;
    runCommandGroup(commandGroup.id);
  };

  const stop = (e: SyntheticEvent) => {
    // Prevent the folder from collapsing when clicking the stop button
    e.stopPropagation();

    if (isReorderingGroups) return;
    stopCommandGroup(commandGroup.id);
  };

  const handleDelete = async () => {
    if (isReorderingGroups) return;
    await deleteCommandGroup(commandGroup.id);
  };

  const handleEdit = () => {
    if (isReorderingGroups) return;
    startEditingCommandGroup(commandGroup);
  };

  return (
    <SidebarGroup
      className="py-0"
      key={commandGroup.id}
      style={style}
      ref={setNodeRef}
    >
      <ContextMenu>
        <ContextMenuTrigger disabled={isReorderingGroups}>
          <SidebarGroupLabel
            asChild
            className="text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground text-sm"
          >
            <button
              className="group flex items-center gap-2 p-2 w-full justify-between"
              style={{ cursor: isReorderingGroups ? "grabbing" : "pointer" }}
              onClick={() => setInternalIsOpen(!internalIsOpen)}
              disabled={isReorderingGroups}
              {...(isReorderingGroups ? { ...attributes, ...listeners } : {})}
            >
              <div className="flex items-center gap-2">
                {isOpen ? <FolderOpen size={16} /> : <Folder size={16} />}
                <p>{commandGroup.name}</p>
              </div>
              {!isReorderingGroups && (
                <div className="flex gap-2 items-center">
                  {someCommandIsRunning && (
                    <span>
                      ({numberOfCommandsRunning}/{commandGroup.commands.length})
                    </span>
                  )}
                  {someCommandIsIdle && (
                    <Play
                      size={16}
                      className="text-muted-foreground cursor-pointer hover:text-primary"
                      onClick={run}
                    />
                  )}
                  {someCommandIsRunning && (
                    <div className="group/command p-0 m-0">
                      <Square
                        size={16}
                        className="text-muted-foreground cursor-pointer hover:text-primary"
                        onClick={stop}
                      />
                    </div>
                  )}
                </div>
              )}
            </button>
          </SidebarGroupLabel>
        </ContextMenuTrigger>
        <ContextMenuContent>
          <ContextMenuItem onClick={handleEdit}>Edit</ContextMenuItem>
          <ContextMenuItem onClick={handleDelete}>Delete</ContextMenuItem>
        </ContextMenuContent>
      </ContextMenu>
      {isOpen && (
        <SidebarGroupContent>
          <SidebarMenu>
            {commandGroup.commands.map((command) => (
              <SidebarMenuItem key={command.id}>
                <CommandMenuItem command={command} />
              </SidebarMenuItem>
            ))}
          </SidebarMenu>
        </SidebarGroupContent>
      )}
    </SidebarGroup>
  );
};
