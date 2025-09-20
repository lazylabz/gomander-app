import { DndContext, type DragEndEvent } from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { ChevronDown, ChevronRight } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { CommandMenuItem } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/CommandMenuItem/CommandMenuItem.tsx";
import { CreateCommandModal } from "@/components/modals/Command/CreateCommandModal.tsx";
import { ALL_COMMANDS_SECTION_OPEN } from "@/constants/localStorage.ts";
import type { Command } from "@/contracts/types.ts";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/design-system/components/ui/collapsible.tsx";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "@/design-system/components/ui/context-menu.tsx";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuItem,
} from "@/design-system/components/ui/sidebar.tsx";
import { parseError } from "@/helpers/errorHelpers.ts";
import { useLocalStorageState } from "@/hooks/useLocalStorageState.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { useCommandStore } from "@/store/commandStore.ts";
import { reorderCommands } from "@/useCases/command/reorderCommands.ts";

export const AllCommandsSection = () => {
  const { t } = useTranslation();
  const commands = useCommandStore((state) => state.commands);
  const setCommands = useCommandStore((state) => state.setCommands);

  const [modalOpen, setModalOpen] = useState(false);

  const [sectionOpen, setSectionOpen] = useLocalStorageState(
    ALL_COMMANDS_SECTION_OPEN,
    false,
  );

  const openCreateCommandModal = () => {
    setModalOpen(true);
  };

  const handleSaveReorderedCommands = async (reorderedCommands: Command[]) => {
    const newOrder = reorderedCommands.map((command) => command.id);
    setCommands(reorderedCommands);
    try {
      await reorderCommands(newOrder);
    } catch (e) {
      toast.error(parseError(e, "Failed to reorder commands"));
    } finally {
      fetchCommands();
    }
  };

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;

    const shouldReorder = active.id && over?.id && active.id !== over.id;
    if (!shouldReorder) {
      return;
    }

    // 1. Find old and new indexes of the dragged command
    const oldIndex = commands.findIndex(
      (command) => command.id === active.id.toString(),
    );
    const newIndex = commands.findIndex(
      (command) => command.id === over.id.toString(),
    );
    if (oldIndex === -1 || newIndex === -1) {
      toast.error("Invalid drag operation: command not found");
      return;
    }

    // 2. Reorder the commands array
    const reorderedCommands = arrayMove(commands, oldIndex, newIndex);

    // 3. Persist the new order
    await handleSaveReorderedCommands(reorderedCommands);
  };

  return (
    <>
      <CreateCommandModal open={modalOpen} setOpen={setModalOpen} />
      <Collapsible
        className="group/collapsible"
        open={sectionOpen}
        onOpenChange={setSectionOpen}
      >
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
                  <p>{t('sidebar.commands.title')}</p>
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
