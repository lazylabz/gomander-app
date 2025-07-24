import { DndContext, type DragEndEvent } from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { ChevronRight, Folder, GripVertical, Play, Square } from "lucide-react";
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
import type { CommandGroup } from "@/contracts/types.ts";
import { useCommandGroupStore } from "@/store/commandGroupStore.ts";
import { useCommandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { runCommandGroup } from "@/useCases/commandGroup/runCommandGroup";
import { saveCommandGroups } from "@/useCases/commandGroup/saveCommandGroups.ts";
import { stopCommandGroup } from "@/useCases/commandGroup/stopCommandGroup.ts";

export const CommandGroupSection = ({
  commandGroup,
}: {
  commandGroup: CommandGroup;
}) => {
  const commandGroups = useCommandGroupStore((state) => state.commandGroups);

  const commands = useCommandStore((state) => state.commands);
  const commandsStatus = useCommandStore((state) => state.commandsStatus);

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

  const editCommandGroupCommands = async (commandGroup: CommandGroup) => {
    await saveCommandGroups(
      commandGroups.map((cg) =>
        cg.id === commandGroup.id ? commandGroup : cg,
      ),
    );
  };

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;

    if (active.id && over?.id && active.id !== over.id) {
      const oldIndex = commandGroup.commands.findIndex(
        (cmdId) => cmdId === active.id,
      );
      const newIndex = commandGroup.commands.findIndex(
        (cmdId) => cmdId === over.id,
      );
      const newCommandsGroups = arrayMove(
        commandGroup.commands,
        oldIndex,
        newIndex,
      );
      await editCommandGroupCommands({
        ...commandGroup,
        commands: newCommandsGroups,
      });
    }
  };

  return (
    <Collapsible
      key={commandGroup.id}
      className="group/collapsible"
      style={style}
      ref={setNodeRef}
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
                  <div
                    {...attributes}
                    {...listeners}
                    className="cursor-grab active:cursor-grabbing pr-0.5 rounded hover:bg-sidebar-accent/50 group-data-[state=open]:hidden transition-transform"
                  >
                    <GripVertical size={14} className="text-muted-foreground" />
                  </div>
                  <ChevronRight
                    size={16}
                    className="transition-transform group-data-[state=open]:rotate-90 group-data-[state=closed]:hidden"
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
              <DndContext onDragEnd={handleDragEnd}>
                <SortableContext
                  strategy={verticalListSortingStrategy}
                  items={commandGroup.commands}
                >
                  {commandGroup.commands
                    .map((cid) => Object(commands[cid]))
                    .filter(Boolean)
                    .map((command) => (
                      <SidebarMenuItem key={command.id}>
                        <CommandMenuItem draggable command={command} />
                      </SidebarMenuItem>
                    ))}
                </SortableContext>
              </DndContext>
            </SidebarMenu>
          </SidebarGroupContent>
        </CollapsibleContent>
      </SidebarGroup>
    </Collapsible>
  );
};
