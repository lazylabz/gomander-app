import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { ChevronRight, Folder, Play, Square } from "lucide-react";
import type { SyntheticEvent } from "react";

import { CommandMenuItem } from "@/components/layout/AppSidebar/components/CommandMenuItem/CommandMenuItem.tsx";
import { useSidebarContext } from "@/components/layout/AppSidebar/contexts/sidebarContext.tsx";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible.tsx";
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
import { CommandStatus, useDataContext } from "@/contexts/DataContext.tsx";
import type { CommandGroup } from "@/types/contracts.ts";

export const CommandGroupSection = ({
  commandGroup,
}: {
  commandGroup: CommandGroup;
}) => {
  const {
    commands,
    commandsStatus,
    runCommandGroup,
    stopCommandGroup,
    saveCommandGroups,
    commandGroups,
  } = useDataContext();

  const { editCommandGroup } = useSidebarContext();

  const { attributes, listeners, setNodeRef, transform, transition } =
    useSortable({ id: commandGroup.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  const numberOfCommandsRunning = commandGroup.commands.filter(
    (commandId) => commandsStatus[commandId] === CommandStatus.RUNNING,
  ).length;

  const someCommandIsRunning = numberOfCommandsRunning > 0;

  const someCommandIsIdle = commandGroup.commands.some(
    (commandId) => commandsStatus[commandId] === CommandStatus.IDLE,
  );

  const run = (e: SyntheticEvent) => {
    // Prevent the folder from collapsing when clicking the play button
    e.stopPropagation();
    runCommandGroup(commandGroup.id);
  };

  const stop = (e: SyntheticEvent) => {
    // Prevent the folder from collapsing when clicking the stop button
    e.stopPropagation();

    stopCommandGroup(commandGroup.id);
  };

  const deleteCommandGroup = async () => {
    await saveCommandGroups(
      commandGroups.filter((cg) => cg.id !== commandGroup.id),
    );
  };

  const handleEdit = () => {
    editCommandGroup(commandGroup);
  };

  return (
    <Collapsible
      key={commandGroup.id}
      className="group/collapsible"
      ref={setNodeRef}
      {...attributes}
      {...listeners}
      style={style}
    >
      <SidebarGroup className="py-0">
        <ContextMenu>
          <ContextMenuTrigger>
            <SidebarGroupLabel
              asChild
              className="group/label text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground text-sm"
            >
              <CollapsibleTrigger className="group flex items-center gap-2 p-2 w-full justify-between">
                <div className="flex items-center gap-2">
                  <ChevronRight
                    size={16}
                    className="transition-transform group-data-[state=open]:rotate-90"
                  />
                  <Folder size={16} />
                  <p>{commandGroup.name}</p>
                </div>
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
              </CollapsibleTrigger>
            </SidebarGroupLabel>
          </ContextMenuTrigger>
          <ContextMenuContent>
            <ContextMenuItem onClick={handleEdit}>Edit</ContextMenuItem>
            <ContextMenuItem onClick={deleteCommandGroup}>
              Delete
            </ContextMenuItem>
          </ContextMenuContent>
        </ContextMenu>
        <CollapsibleContent className="pl-4">
          <SidebarGroupContent>
            <SidebarMenu>
              {Object.values(commands)
                .filter((c) => commandGroup.commands.includes(c.id))
                .map((command) => (
                  <SidebarMenuItem key={command.id}>
                    <CommandMenuItem command={command} />
                  </SidebarMenuItem>
                ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </CollapsibleContent>
      </SidebarGroup>
    </Collapsible>
  );
};
