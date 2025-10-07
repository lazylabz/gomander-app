import { ChevronDownIcon, Import, Plus } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";

import { CreateProjectModal } from "@/components/modals/Project/CreateProjectModal.tsx";
import { DeleteProjectModal } from "@/components/modals/Project/DeleteProjectModal.tsx";
import { ImportProjectModal } from "@/components/modals/Project/ImportProjectModal.tsx";
import { dataService } from "@/contracts/service.ts";
import type { ProjectExport } from "@/contracts/types.ts";
import { Button } from "@/design-system/components/ui/button.tsx";
import { ButtonGroup } from "@/design-system/components/ui/button-group.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/design-system/components/ui/dropdown-menu.tsx";
import { parseError } from "@/helpers/errorHelpers.ts";
import { ScreenRoutes } from "@/routes.ts";
import { ProjectCard } from "@/screens/ProjectSelectionScreen/components/ProjectCard.tsx";
import { useGetAvailableProjects } from "@/screens/ProjectSelectionScreen/hooks/useGetAvailableProjects.ts";
import { useProjectStore } from "@/store/projectStore.ts";
import { deleteProject } from "@/useCases/project/deleteProject.ts";

export const ProjectSelectionScreen = () => {
  const [projectIdBeingDeleted, setProjectIdBeingDeleted] = useState<
    string | null
  >(null);
  const project = useProjectStore((state) => state.projectInfo);

  const navigate = useNavigate();

  const [projectBeingImported, setProjectBeingImported] =
    useState<ProjectExport | null>(null);

  const [createProjectModalOpen, setCreateProjectModalOpen] = useState(false);

  const { availableProjects, fetchAvailableProjects } =
    useGetAvailableProjects();

  const openCreateProjectModal = () => {
    setCreateProjectModalOpen(true);
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
    try {
      const projectToImport = await dataService.getProjectToImport();
      setProjectBeingImported(projectToImport);
    } catch (e) {
      toast.error(parseError(e, "Failed to select project"));
    }
  };

  const handleImportProjectFromPackageJson = async () => {
    try {
      const projectToImport =
        await dataService.getProjectToImportFromPackageJson();
      setProjectBeingImported(projectToImport);
    } catch (e) {
      toast.error(parseError(e, "Failed to select project"));
    }
  };

  useEffect(() => {
    fetchAvailableProjects();
  }, [fetchAvailableProjects]);

  useEffect(() => {
    if (project) {
      navigate(ScreenRoutes.Logs);
    }
  }, [navigate, project]);

  const hasProjects = availableProjects.length > 0;

  return (
    <>
      <CreateProjectModal
        onSuccess={fetchAvailableProjects}
        open={createProjectModalOpen}
        setOpen={setCreateProjectModalOpen}
      />
      <ImportProjectModal
        open={!!projectBeingImported}
        onSuccess={fetchAvailableProjects}
        onClose={() => setProjectBeingImported(null)}
        project={projectBeingImported}
      />
      <DeleteProjectModal
        open={!!projectIdBeingDeleted}
        onConfirm={confirmDeleteProject}
        onClose={cancelDeleteProject}
      />
      <div className="w-full h-full flex flex-col items-center justify-center gap-10">
        <h1 className="text-3xl">
          {hasProjects ? "Open project" : "Welcome to Gomander!"}
        </h1>
        {hasProjects && (
          <div className="flex flex-col items-center justify-center gap-2">
            {availableProjects.map((p) => (
              <ProjectCard
                key={p.id}
                project={p}
                handleDeleteProject={handleDeleteProject(p.id)}
              />
            ))}
          </div>
        )}
        {!hasProjects && (
          <p>You don't have projects yet. Create or import one.</p>
        )}
        <div className="flex flex-row items-center gap-2 justify-center">
          <Button onClick={openCreateProjectModal} variant="ghost">
            <Plus /> Create a new project
          </Button>
          <ButtonGroup>
            <Button variant="ghost" onClick={handleImportProject}>
              <Import />
              Import an existing project
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" aria-label="More Options">
                  <ChevronDownIcon />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-52">
                <DropdownMenuItem onClick={handleImportProjectFromPackageJson}>
                  From a package.json
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </ButtonGroup>
        </div>
      </div>
    </>
  );
};
