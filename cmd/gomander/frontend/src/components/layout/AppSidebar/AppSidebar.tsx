import { DndContext, type DragEndEvent } from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { ChevronDown, Settings } from "lucide-react";
import { useState } from "react";

import { AllCommandsSection } from "@/components/layout/AppSidebar/components/AllCommandsSection/AllCommandsSection.tsx";
import { CommandGroupSection } from "@/components/layout/AppSidebar/components/CommandGroupSection/CommandGroupSection.tsx";
import { CreateMenu } from "@/components/layout/AppSidebar/components/CreateMenu/CreateMenu.tsx";
import { sidebarContext } from "@/components/layout/AppSidebar/contexts/sidebarContext.tsx";
import { EditCommandModal } from "@/components/modals/Command/EditCommandModal.tsx";
import { EditCommandGroupModal } from "@/components/modals/CommandGroup/EditCommandGroupModal.tsx";
import { EditOpenedProjectModal } from "@/components/modals/Project/EditOpenedProjectModal.tsx";
import { SettingsModal } from "@/components/modals/SettingsModal.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
} from "@/components/ui/sidebar.tsx";
import type { Command, CommandGroup, ProjectInfo } from "@/contracts/types.ts";
import { fetchProject } from "@/queries/fetchProject.ts";
import { useCommandGroupStore } from "@/store/commandGroupStore.ts";
import { useProjectStore } from "@/store/projectStore.ts";
import { saveCommandGroups } from "@/useCases/commandGroup/saveCommandGroups.ts";
import { closeProject } from "@/useCases/project/closeProject.ts";

export const AppSidebar = () => {
  const commandGroups = useCommandGroupStore((state) => state.commandGroups);
  const project = useProjectStore((state) => state.projectInfo);

  const [editingProject, setEditingProject] = useState<ProjectInfo | null>(
    null,
  );
  const [editingCommand, setEditingCommand] = useState<Command | null>(null);
  const [editingCommandGroup, setEditingCommandGroup] =
    useState<CommandGroup | null>(null);
  const [settingsModalOpen, setSettingsModalOpen] = useState(false);

  const openEditProjectModal = () => {
    setEditingProject(project);
  };

  const closeEditProjectModal = () => {
    setEditingProject(null);
  };

  const closeEditCommandModal = () => {
    setEditingCommand(null);
  };

  const closeEditCommandGroupModal = () => {
    setEditingCommandGroup(null);
  };

  const openSettingsModal = () => {
    setSettingsModalOpen(true);
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
    <sidebarContext.Provider
      value={{
        startEditingCommand: (command: Command) => setEditingCommand(command),
        startEditingCommandGroup: (commandGroup: CommandGroup) =>
          setEditingCommandGroup(commandGroup),
      }}
    >
      <EditOpenedProjectModal
        open={!!editingProject}
        onClose={closeEditProjectModal}
        onSuccess={fetchProject}
        project={editingProject}
      />
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
          <div className="flex items-center ml-2 gap-2">
            <p className="text-xl font-extralight">G.</p>
            <DropdownMenu>
              <DropdownMenuTrigger className="flex gap-1 items-center hover:bg-sidebar-foreground/8 p-1 px-1 pr-2 rounded-md">
                <h1 className="text-xl font-semibold pl-2">{project?.name}</h1>
                <ChevronDown className="mt-1" size={20} />
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <DropdownMenuItem onClick={openEditProjectModal}>
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem onClick={closeProject}>
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
