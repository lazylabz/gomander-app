import { Settings } from "lucide-react";
import { useState } from "react";

import { CommandGroupSection } from "@/components/layout/AppSidebar/components/CommandGroupSection/CommandGroupSection.tsx";
import { CreateMenu } from "@/components/layout/AppSidebar/components/CreateMenu/CreateMenu.tsx";
import { MyCommandsSection } from "@/components/layout/AppSidebar/components/MyCommandsSection/MyCommandsSection.tsx";
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
import { useDataContext } from "@/contexts/DataContext.tsx";
import type { Command, CommandGroup } from "@/types/contracts.ts";

export const AppSidebar = () => {
  const { commandGroups } = useDataContext();

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
          <MyCommandsSection />
          {commandGroups.map((cg) => (
            <CommandGroupSection commandGroup={cg} key={cg.id} />
          ))}
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
