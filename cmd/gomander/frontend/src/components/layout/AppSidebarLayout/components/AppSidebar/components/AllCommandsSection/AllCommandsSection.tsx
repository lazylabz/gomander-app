import { DndContext, type DragEndEvent } from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { ChevronDown, ChevronRight } from "lucide-react";
import { useState } from "react";

import { CommandMenuItem } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/CommandMenuItem/CommandMenuItem.tsx";
import { CreateCommandModal } from "@/components/modals/Command/CreateCommandModal.tsx";
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
import { useCommandStore } from "@/store/commandStore.ts";
import { reorderCommands } from "@/useCases/command/reorderCommands.ts";

export const AllCommandsSection = () => {
  const commands = useCommandStore((state) => state.commands);

  const [modalOpen, setModalOpen] = useState(false);

  const openCreateCommandModal = () => {
    setModalOpen(true);
  };

  const handleSaveReorderedCommands = async (newOrder: string[]) => {
    await reorderCommands(newOrder);
  };

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;

    const shouldReorder = active.id && over?.id && active.id !== over.id;
    if (!shouldReorder) {
      return;
    }

    const commandsIds = commands.map((command) => command.id);
    // 1. Find old and new indexes of the dragged command
    const oldIndex = commandsIds.indexOf(active.id.toString());
    const newIndex = commandsIds.indexOf(over.id.toString());

    // 2. Reorder the commands array
    const newOrder = arrayMove(commandsIds, oldIndex, newIndex);

    // 3. Persist the new order
    await handleSaveReorderedCommands(newOrder);
  };

  return (
    <>
      <CreateCommandModal open={modalOpen} setOpen={setModalOpen} />
      <Collapsible className="group/collapsible">
        <SidebarGroup className="py-0">
          <ContextMenu>
            <ContextMenuTrigger>
              <SidebarGroupLabel
                asChild
                className="group/label text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground text-sm"
              >
                <CollapsibleTrigger className="group flex items-center gap-2 p-2 pl-1 w-full">
                  <ChevronDown className="hidden group-data-[state=open]:block" />
                  <ChevronRight className="block group-data-[state=open]:hidden" />
                  <p>All commands</p>
                </CollapsibleTrigger>
              </SidebarGroupLabel>
            </ContextMenuTrigger>
            <ContextMenuContent>
              <ContextMenuItem onClick={openCreateCommandModal}>
                Add command
              </ContextMenuItem>
            </ContextMenuContent>
          </ContextMenu>
          <CollapsibleContent>
            <SidebarGroupContent>
              <SidebarMenu>
                <DndContext onDragEnd={handleDragEnd}>
                  <SortableContext
                    strategy={verticalListSortingStrategy}
                    items={commands}
                  >
                    {commands.map((command) => (
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
    </>
  );
};
