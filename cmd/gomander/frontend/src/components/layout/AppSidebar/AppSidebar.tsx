import { DndContext, type DragEndEvent } from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { Settings } from "lucide-react";
import { useState } from "react";

import { AllCommandsSection } from "@/components/layout/AppSidebar/components/AllCommandsSection/AllCommandsSection.tsx";
import { CommandGroupSection } from "@/components/layout/AppSidebar/components/CommandGroupSection/CommandGroupSection.tsx";
import { CreateMenu } from "@/components/layout/AppSidebar/components/CreateMenu/CreateMenu.tsx";
import { sidebarContext } from "@/components/layout/AppSidebar/contexts/sidebarContext.tsx";
import { EditCommandModal } from "@/components/modals/Command/EditCommandModal.tsx";
import { EditCommandGroupModal } from "@/components/modals/CommandGroup/EditCommandGroupModal.tsx";
import { SettingsModal } from "@/components/modals/SettingsModal.tsx";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
} from "@/components/ui/sidebar.tsx";
import type { Command, CommandGroup } from "@/contracts/types.ts";
import { useCommandGroupStore } from "@/store/commandGroupStore.ts";
import { saveCommandGroups } from "@/useCases/commandGroup/saveCommandGroups.ts";

export const AppSidebar = () => {
  const commandGroups = useCommandGroupStore((state) => state.commandGroups);

  const [editingCommand, setEditingCommand] = useState<Command | null>(null);
  const [editingCommandGroup, setEditingCommandGroup] =
    useState<CommandGroup | null>(null);
  const [settingsModalOpen, setSettingsModalOpen] = useState(false);

  const closeEditCommandModal = () => {
    setEditingCommand(null);
  };

  const closeEditCommandGroupModal = () => {
    setEditingCommandGroup(null);
  };

  const openSettingsModal = () => {
    setSettingsModalOpen(true);
  };

  const value = {
    editCommand: (command: Command) => setEditingCommand(command),
    editCommandGroup: (commandGroup: CommandGroup) =>
      setEditingCommandGroup(commandGroup),
  };

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;

    if (active.id && over?.id && active.id !== over.id) {
      const oldIndex = commandGroups.findIndex((cg) => cg.id === active.id);
      const newIndex = commandGroups.findIndex((cg) => cg.id === over.id);
      const newCommandGroups = arrayMove(commandGroups, oldIndex, newIndex);
      await saveCommandGroups(newCommandGroups);
    }
  };

  return (
    <sidebarContext.Provider value={value}>
      <EditCommandModal
        command={editingCommand}
        open={!!editingCommand}
        setOpen={closeEditCommandModal}
      />
      <EditCommandGroupModal
        commandGroup={editingCommandGroup}
        open={!!editingCommandGroup}
        setOpen={closeEditCommandGroupModal}
      />
      <SettingsModal open={settingsModalOpen} setOpen={setSettingsModalOpen} />
      <Sidebar collapsible="icon">
        <SidebarHeader className="flex flex-row items-center justify-between p-2">
          <h1 className="text-xl font-semibold pl-2">Gomander</h1>
          <CreateMenu />
        </SidebarHeader>
        <SidebarContent className="gap-0">
          <AllCommandsSection />
          <h3 className="text-sm pl-4 mt-4 mb-1 text-muted-foreground">
            Command groups
          </h3>
          <DndContext onDragEnd={handleDragEnd}>
            <SortableContext
              items={commandGroups.map((cg) => cg.id)}
              strategy={verticalListSortingStrategy}
            >
              {commandGroups.map((cg) => (
                <CommandGroupSection commandGroup={cg} key={cg.id} />
              ))}
            </SortableContext>
          </DndContext>
        </SidebarContent>
        <SidebarFooter>
          <Settings
            onClick={openSettingsModal}
            size={20}
            className="text-muted-foreground cursor-pointer hover:text-primary"
          />
        </SidebarFooter>
      </Sidebar>
    </sidebarContext.Provider>
  );
};
