import { ChevronRight, Folder, Play } from "lucide-react";

import { CreateMenu } from "@/components/layout/AppSidebar/components/CreateMenu/CreateMenu.tsx";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible.tsx";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
} from "@/components/ui/context-menu";
import { ContextMenuTrigger } from "@/components/ui/context-menu.tsx";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar.tsx";
import { useDataContext } from "@/contexts/DataContext.tsx";
import type { Command } from "@/types/contracts.ts";
import { EditCommandModal } from "@/components/modals/EditCommandModal.tsx";
import { createContext, useContext, useState } from "react";

const editContext = createContext<{
  editCommand: (command: Command) => void;
}>(
  {} as {
    editCommand: (command: Command) => void;
  },
);

const CommandMenuItem = ({ command }: { command: Command }) => {
  const { execCommand, setActiveCommandId, activeCommandId, deleteCommand } =
    useDataContext();

  const { editCommand } = useContext(editContext);

  const handleCommandClick = async () => {
    setActiveCommandId(command.id);
    await execCommand(command.id);
  };

  const handleDeleteCommand = async () => {
    await deleteCommand(command.id);
    setActiveCommandId(null); // Reset active command after deletion
  };

  const handleEditCommand = () => {
    editCommand(command);
  };

  const onCommandSectionClick = () => {
    setActiveCommandId(command.id);
  };

  return (
    <ContextMenu>
      <ContextMenuTrigger>
        <SidebarMenuButton asChild isActive={activeCommandId === command.id}>
          <div
            onClick={onCommandSectionClick}
            className="flex flex-row justify-between items-center w-full"
          >
            {command.name}
            <Play
              className="text-muted-foreground cursor-pointer hover:text-primary"
              onClick={handleCommandClick}
            />
          </div>
        </SidebarMenuButton>
      </ContextMenuTrigger>
      <ContextMenuContent>
        <ContextMenuItem onClick={handleEditCommand}>Edit</ContextMenuItem>
        <ContextMenuItem onClick={handleDeleteCommand}>Delete</ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  );
};

export const MyCommandsSection = () => {
  const { commands } = useDataContext();

  return (
    <Collapsible className="group/collapsible">
      <SidebarGroup>
        <SidebarGroupLabel
          asChild
          className="group/label text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground text-sm"
        >
          <CollapsibleTrigger className="group flex items-center gap-2 p-2">
            <ChevronRight className="transition-transform group-data-[state=open]:rotate-90" />
            <Folder />
            <p>Your commands</p>
          </CollapsibleTrigger>
        </SidebarGroupLabel>
        <CollapsibleContent className="pl-4">
          <SidebarGroupContent>
            <SidebarMenu>
              {Object.values(commands).map((command) => (
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

export const AppSidebar = () => {
  const [editingCommand, setEdittingCommand] = useState<Command | null>(null);

  const closeModal = () => {
    setEdittingCommand(null);
  };

  const value = {
    editCommand: (command: Command) => setEdittingCommand(command),
  };

  return (
    <editContext.Provider value={value}>
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
    </editContext.Provider>
  );
};
