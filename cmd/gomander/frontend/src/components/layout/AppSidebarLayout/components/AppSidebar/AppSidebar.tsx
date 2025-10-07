import { ChevronDown, Settings } from "lucide-react";
import { useState } from "react";
import { useNavigate } from "react-router";

import { AllCommandsSection } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/AllCommandsSection/AllCommandsSection.tsx";
import { CommandGroupsSection } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/CommandGroupsSection/CommandGroupsSection.tsx";
import { CreateMenu } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/CreateMenu/CreateMenu.tsx";
import { VersionSection } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/VersionSection/VersionSection.tsx";
import { sidebarContext } from "@/components/layout/AppSidebarLayout/components/AppSidebar/contexts/sidebarContext.tsx";
import { AboutModal } from "@/components/modals/About/AboutModal.tsx";
import { EditCommandModal } from "@/components/modals/Command/EditCommandModal.tsx";
import type { Command } from "@/contracts/types.ts";
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/design-system/components/ui/avatar.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/design-system/components/ui/dropdown-menu.tsx";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
} from "@/design-system/components/ui/sidebar.tsx";
import { ScreenRoutes } from "@/routes.ts";
import { SettingsTab } from "@/screens/SettingsScreen/contexts/settingsContext.tsx";
import { useProjectStore } from "@/store/projectStore.ts";
import { closeProject } from "@/useCases/project/closeProject.ts";

export const AppSidebar = () => {
  const project = useProjectStore((state) => state.projectInfo);

  const navigate = useNavigate();

  const [editingCommand, setEditingCommand] = useState<Command | null>(null);

  const [aboutModalOpen, setAboutModalOpen] = useState(false);

  const closeEditCommandModal = () => {
    setEditingCommand(null);
  };

  const goToSettings = (tab: SettingsTab) => {
    navigate(ScreenRoutes.Settings, { state: { tab } });
  };

  const handleCloseProject = async () => {
    await closeProject();
    navigate(ScreenRoutes.ProjectSelection);
  };

  return (
    <sidebarContext.Provider
      value={{
        startEditingCommand: (command: Command) => setEditingCommand(command),
      }}
    >
      <EditCommandModal
        command={editingCommand}
        open={!!editingCommand}
        setOpen={closeEditCommandModal}
      />
      <AboutModal open={aboutModalOpen} setOpen={setAboutModalOpen} />
      <Sidebar collapsible="offcanvas">
        <SidebarHeader className="flex flex-row items-center justify-between p-2">
          <div className="flex items-center ml-2 gap-1">
            <Avatar className="size-9 rounded-lg mb-1">
              <AvatarImage src="/sidebar-logo.png" />
              <AvatarFallback className="text-xl font-extralight">
                G.
              </AvatarFallback>
            </Avatar>
            <DropdownMenu>
              <DropdownMenuTrigger className="flex gap-1 items-center hover:bg-sidebar-foreground/8 p-1 px-2 rounded-md">
                <h1 className="text-xl font-semibold pl-1">{project?.name}</h1>
                <ChevronDown className="mt-1" size={20} />
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <DropdownMenuItem
                  onClick={() => goToSettings(SettingsTab.Project)}
                >
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem onClick={handleCloseProject}>
                  Close
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
          <CreateMenu />
        </SidebarHeader>
        <SidebarContent className="gap-1">
          <AllCommandsSection />
          <CommandGroupsSection />
        </SidebarContent>
        <SidebarFooter className="flex flex-row items-center justify-between p-2">
          <Settings
            onClick={() => goToSettings(SettingsTab.User)}
            size={20}
            className="text-muted-foreground cursor-pointer hover:text-primary"
          />
          <VersionSection openAboutModal={() => setAboutModalOpen(true)} />
        </SidebarFooter>
      </Sidebar>
    </sidebarContext.Provider>
  );
};
