import { DndContext, type DragEndEvent } from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { ChevronDown, Settings } from "lucide-react";
import { useState } from "react";
import { useNavigate } from "react-router";

import { AllCommandsSection } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/AllCommandsSection/AllCommandsSection.tsx";
import { CommandGroupSection } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/CommandGroupSection/CommandGroupSection.tsx";
import { CreateMenu } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/CreateMenu/CreateMenu.tsx";
import { VersionSection } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/VersionSection/VersionSection.tsx";
import { sidebarContext } from "@/components/layout/AppSidebarLayout/components/AppSidebar/contexts/sidebarContext.tsx";
import { AboutModal } from "@/components/modals/About/AboutModal.tsx";
import { EditCommandModal } from "@/components/modals/Command/EditCommandModal.tsx";
import { EditCommandGroupModal } from "@/components/modals/CommandGroup/EditCommandGroupModal.tsx";
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/components/ui/avatar.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
} from "@/components/ui/sidebar.tsx";
import type { Command, CommandGroup } from "@/contracts/types.ts";
import { ScreenRoutes } from "@/routes.ts";
import { SettingsTab } from "@/screens/SettingsScreen/contexts/settingsContext.tsx";
import { useCommandGroupStore } from "@/store/commandGroupStore.ts";
import { useProjectStore } from "@/store/projectStore.ts";
import { saveCommandGroups } from "@/useCases/commandGroup/saveCommandGroups.ts";
import { closeProject } from "@/useCases/project/closeProject.ts";

export const AppSidebar = () => {
  const commandGroups = useCommandGroupStore((state) => state.commandGroups);
  const project = useProjectStore((state) => state.projectInfo);

  const navigate = useNavigate();

  const [editingCommand, setEditingCommand] = useState<Command | null>(null);
  const [editingCommandGroup, setEditingCommandGroup] =
    useState<CommandGroup | null>(null);

  const [aboutModalOpen, setAboutModalOpen] = useState(false);

  const closeEditCommandModal = () => {
    setEditingCommand(null);
  };

  const closeEditCommandGroupModal = () => {
    setEditingCommandGroup(null);
  };

  const goToSettings = (tab: SettingsTab) => {
    navigate(ScreenRoutes.Settings, { state: { tab } });
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

  const handleCloseProject = async () => {
    await closeProject();
    navigate(ScreenRoutes.ProjectSelection);
  };

  return (
    <sidebarContext.Provider
      value={{
        startEditingCommand: (command: Command) => setEditingCommand(command),
        startEditingCommandGroup: (commandGroup: CommandGroup) =>
          setEditingCommandGroup(commandGroup),
      }}
    >
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
      <AboutModal open={aboutModalOpen} setOpen={setAboutModalOpen} />
      <Sidebar collapsible="icon">
        <SidebarHeader className="flex flex-row items-center justify-between p-2">
          <div className="flex items-center ml-2 gap-1">
            <Avatar className="size-9 rounded-lg">
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
