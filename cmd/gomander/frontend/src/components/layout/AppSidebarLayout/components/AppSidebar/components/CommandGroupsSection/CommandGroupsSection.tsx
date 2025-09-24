import { DndContext, type DragEndEvent } from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { ArrowUpDown } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

import { CommandGroupSection } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/CommandGroupSection/CommandGroupSection.tsx";
import { CreateCommandGroupModal } from "@/components/modals/CommandGroup/CreateCommandGroupModal.tsx";
import { EditCommandGroupModal } from "@/components/modals/CommandGroup/EditCommandGroupModal.tsx";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "@/components/ui/context-menu.tsx";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip.tsx";
import type { CommandGroup } from "@/contracts/types.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { cn } from "@/lib/utils.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { useCommandGroupStore } from "@/store/commandGroupStore.ts";
import { reorderCommandGroups } from "@/useCases/commandGroup/reorderCommandGroups.ts";

export const CommandGroupsSection = () => {
  const commandGroups = useCommandGroupStore((state) => state.commandGroups);
  const setCommandGroups = useCommandGroupStore(
    (state) => state.setCommandGroups,
  );

  const [isReorderingGroups, setIsReorderingGroups] = useState(false);
  const [editingCommandGroup, setEditingCommandGroup] =
    useState<CommandGroup | null>(null);
  const [createCommandGroupModalOpen, setCreateCommandGroupModalOpen] =
    useState(false);

  const openCreateCommandGroupModal = () => {
    setCreateCommandGroupModalOpen(true);
  };

  const startEditingCommandGroup = (commandGroup: CommandGroup) => {
    setEditingCommandGroup(commandGroup);
  };

  const closeEditCommandGroupModal = () => {
    setEditingCommandGroup(null);
  };

  const handleCommandGroupDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;

    if (active.id && over?.id && active.id !== over.id) {
      const oldIndex = commandGroups.findIndex((cg) => cg.id === active.id);
      const newIndex = commandGroups.findIndex((cg) => cg.id === over.id);
      const newCommandGroups = arrayMove(commandGroups, oldIndex, newIndex);
      setCommandGroups(newCommandGroups);
    }
  };

  const handleSaveCommandGroupsOrder = async () => {
    const reorderedCommandGroupIds = commandGroups.map((cg) => cg.id);
    try {
      await reorderCommandGroups(reorderedCommandGroupIds);
      toast.success("Reordered command groups successfully");
    } catch (e) {
      toast.error(parseError(e, "Failed to reorder command groups"));
    } finally {
      fetchCommandGroups();
    }
  };

  const toggleReorderingMode = () => {
    setIsReorderingGroups((wasReordering) => {
      if (wasReordering) {
        handleSaveCommandGroupsOrder();
      }
      return !wasReordering;
    });
  };

  return (
    <>
      <CreateCommandGroupModal
        open={createCommandGroupModalOpen}
        setOpen={setCreateCommandGroupModalOpen}
      />
      <EditCommandGroupModal
        commandGroup={editingCommandGroup}
        open={!!editingCommandGroup}
        setOpen={closeEditCommandGroupModal}
      />
      <ContextMenu>
        <ContextMenuTrigger>
          <div className="flex items-center pl-4 pr-2 mt-2 mb-1 gap-2">
            <h3 className="text-sm text-muted-foreground">Command groups</h3>
            <Tooltip delayDuration={1000}>
              <TooltipContent>
                {isReorderingGroups ? "Apply reordering" : "Start reordering"}
              </TooltipContent>
              <TooltipTrigger
                onClick={toggleReorderingMode}
                className={cn(
                  "p-1 rounded hover:bg-sidebar-accent transition-colors border border-transparent",
                  isReorderingGroups
                    ? "text-primary bg-sidebar-accent border-muted-foreground"
                    : "text-muted-foreground hover:text-primary",
                )}
              >
                <ArrowUpDown size={14} />
              </TooltipTrigger>
            </Tooltip>
          </div>
        </ContextMenuTrigger>
        <ContextMenuContent>
          <ContextMenuItem onClick={openCreateCommandGroupModal}>
            Add command group
          </ContextMenuItem>
        </ContextMenuContent>
      </ContextMenu>
      <DndContext onDragEnd={handleCommandGroupDragEnd}>
        <SortableContext
          items={commandGroups.map((cg) => cg.id)}
          strategy={verticalListSortingStrategy}
        >
          {commandGroups.map((cg) => (
            <CommandGroupSection
              isReorderingGroups={isReorderingGroups}
              startEditingCommandGroup={startEditingCommandGroup}
              commandGroup={cg}
              key={cg.id}
            />
          ))}
        </SortableContext>
      </DndContext>
    </>
  );
};
