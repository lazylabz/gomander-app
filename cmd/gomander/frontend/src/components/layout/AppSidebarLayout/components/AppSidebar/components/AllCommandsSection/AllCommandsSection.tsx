import { Folder, FolderOpen } from "lucide-react";
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

export const AllCommandsSection = () => {
  const commands = useCommandStore((state) => state.commands);

  const [modalOpen, setModalOpen] = useState(false);

  const openCreateCommandModal = () => {
    setModalOpen(true);
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
                <CollapsibleTrigger className="group flex items-center gap-2 p-2 w-full">
                  <FolderOpen className="hidden group-data-[state=open]:block" />
                  <Folder className="block group-data-[state=open]:hidden" />
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
    </>
  );
};
