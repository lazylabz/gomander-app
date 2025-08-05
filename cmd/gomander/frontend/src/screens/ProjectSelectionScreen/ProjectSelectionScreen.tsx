import { Import, Plus } from "lucide-react";
import { useEffect, useState } from "react";

import { CreateProjectModal } from "@/components/modals/Project/CreateProjectModal.tsx";
import { DeleteProjectModal } from "@/components/modals/Project/DeleteProjectModal.tsx";
import { Button } from "@/components/ui/button.tsx";
import { dataService } from "@/contracts/service.ts";
import type { Project } from "@/contracts/types.ts";
import { ProjectCard } from "@/screens/ProjectSelectionScreen/components/ProjectCard.tsx";
import { deleteProject } from "@/useCases/project/deleteProject.ts";
import { importProject } from "@/useCases/project/importProject.ts";

export const ProjectSelectionScreen = () => {
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

  const handleImportProject = async () => {
    await importProject();
    await fetchAvailableProjects();
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
      <div className="w-full h-full flex flex-col items-center justify-center gap-10">
        <h1 className="text-3xl">Open project </h1>
        {availableProjects.length > 0 && (
          <div className="flex flex-col items-center justify-center gap-2">
            {availableProjects.map((p) => (
              <ProjectCard
                project={p}
                handleDeleteProject={handleDeleteProject(p.id)}
              />
            ))}
          </div>
        )}
        <div className="flex flex-col items-center justify-center gap-2">
          {availableProjects.length === 0 && (
            <p>You don't have projects yet. Create or import one.</p>
          )}
          <div className="flex flex-row items-center gap-2 justify-center">
            <Button onClick={openCreateProjectModal} variant="ghost">
              <Plus /> Create a new project
            </Button>
            <Button onClick={handleImportProject} variant="ghost">
              <Import /> Import an existing project
            </Button>
          </div>
        </div>
        <div></div>
      </div>
    </>
  );
};
