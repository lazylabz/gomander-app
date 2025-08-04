import { EllipsisVertical, Plus } from "lucide-react";
import { useEffect, useState } from "react";

import { CreateProjectModal } from "@/components/modals/Project/CreateProjectModal.tsx";
import { DeleteProjectModal } from "@/components/modals/Project/DeleteProjectModal.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { dataService } from "@/contracts/service.ts";
import type { Project } from "@/contracts/types.ts";
import { fetchProject } from "@/queries/fetchProject.ts";
import { deleteProject } from "@/useCases/project/deleteProject.ts";
import { exportProject } from "@/useCases/project/exportProject.ts";

export const ProjectSelectionScreen = () => {
  const [projectMenuOpen, setProjectMenuOpen] = useState<string | null>(null);
  const [projectIdBeingDeleted, setProjectIdBeingDeleted] = useState<
    string | null
  >(null);

  const [createProjectModalOpen, setCreateProjectModalOpen] = useState(false);
  const [availableProjects, setAvailableProjects] = useState<Project[]>([]);

  const fetchAvailableProjects = async () => {
    const projects = await dataService.getAvailableProjects();
    setAvailableProjects(projects);
  };

  const openCreateProjectModal = () => {
    setCreateProjectModalOpen(true);
  };

  const handleSuccess = async () => {
    await fetchAvailableProjects();
  };

  const handleOpenProject = (projectId: string) => async () => {
    await dataService.openProject(projectId);
    await fetchProject();
  };

  const handleSetProjectMenuOpen = (projectId: string) => (open: boolean) => {
    if (open) {
      setProjectMenuOpen(projectId);
    } else {
      setProjectMenuOpen(null);
    }
  };

  const handleProjectMenuClick = (projectId: string) => () => {
    setProjectMenuOpen(projectId);
  };

  const handleDeleteProject = (projectId: string) => async () => {
    setProjectIdBeingDeleted(projectId);
  };

  const confirmDeleteProject = async () => {
    await deleteProject(projectIdBeingDeleted!);
    setProjectIdBeingDeleted(null);
    await fetchAvailableProjects();
  };

  const cancelDeleteProject = () => {
    setProjectIdBeingDeleted(null);
  };

  useEffect(() => {
    fetchAvailableProjects();
  }, []);

  return (
    <>
      <CreateProjectModal
        onSuccess={handleSuccess}
        open={createProjectModalOpen}
        setOpen={setCreateProjectModalOpen}
      />
      <DeleteProjectModal
        open={!!projectIdBeingDeleted}
        onConfirm={confirmDeleteProject}
        onClose={cancelDeleteProject}
      />
      <div className="w-full h-full flex flex-col items-center justify-center gap-6">
        <h1 className="text-3xl">Open project </h1>
        {availableProjects.length > 0 && (
          <div className="flex flex-col items-center justify-center gap-2">
            {availableProjects.map((p) => (
              <div
                className="relative px-4 py-2 border border-neutral-100 rounded-xl shadow hover:shadow-md transition-all w-80 flex flex-col items-start gap-1"
                key={p.id}
              >
                <DropdownMenu
                  open={projectMenuOpen === p.id}
                  onOpenChange={handleSetProjectMenuOpen(p.id)}
                >
                  <DropdownMenuTrigger className="absolute top-2 right-2 ">
                    <EllipsisVertical
                      onClick={handleProjectMenuClick(p.id)}
                      size={16}
                      className="text-muted-foreground cursor-pointer hover:text-primary"
                    />
                  </DropdownMenuTrigger>
                  <DropdownMenuContent>
                    <DropdownMenuItem onClick={handleDeleteProject(p.id)}>
                      Delete
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => exportProject(p.id)}>
                      Export
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
                <p
                  className="font-semibold hover:underline cursor-pointer"
                  onClick={handleOpenProject(p.id)}
                >
                  {p.name}
                </p>
                <div className="flex text-sm gap-4 text-muted-foreground">
                  <p>{Object.keys(p.commands).length} commands</p>
                  <p>{p.commandGroups.length} command groups</p>
                </div>
              </div>
            ))}
          </div>
        )}
        <div className="flex flex-col items-center justify-center gap-2">
          {availableProjects.length === 0 && <p>You don't have projects yet</p>}
          <Button onClick={openCreateProjectModal} variant="outline">
            <Plus /> Create a new project
          </Button>
        </div>
        <div></div>
      </div>
    </>
  );
};
