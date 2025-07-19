import { useState } from "react";

import { CreateMenu } from "@/components/layout/AppSidebar/components/CreateMenu/CreateMenu.tsx";
import { MyCommandsSection } from "@/components/layout/AppSidebar/components/MyCommandsSection/MyCommandsSection.tsx";
import { sidebarContext } from "@/components/layout/AppSidebar/contexts/sidebarContext.tsx";
import { EditCommandModal } from "@/components/modals/EditCommandModal.tsx";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
} from "@/components/ui/sidebar.tsx";
import type { Command } from "@/types/contracts.ts";

export const AppSidebar = () => {
  const [editingCommand, setEdittingCommand] = useState<Command | null>(null);

  const closeModal = () => {
    setEdittingCommand(null);
  };

  const value = {
    editCommand: (command: Command) => setEdittingCommand(command),
  };

  return (
    <sidebarContext.Provider value={value}>
      <EditCommandModal
        command={editingCommand}
        open={!!editingCommand}
        setOpen={closeModal}
      />
      <Sidebar collapsible="icon">
        <SidebarHeader className="flex flex-row items-center justify-between p-2">
          <h1 className="text-xl font-semibold pl-2">Gomander</h1>
          <CreateMenu />
        </SidebarHeader>
        <SidebarContent>
          <MyCommandsSection />
        </SidebarContent>
        <SidebarFooter />
      </Sidebar>
    </sidebarContext.Provider>
  );
};
