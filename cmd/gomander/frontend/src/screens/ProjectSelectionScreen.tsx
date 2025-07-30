import { Plus } from "lucide-react";
import { useEffect, useState } from "react";

import { CreateProjectModal } from "@/components/modals/Project/CreateProjectModal.tsx";
import { Button } from "@/components/ui/button.tsx";
import { dataService } from "@/contracts/service.ts";
import type { Project } from "@/contracts/types.ts";
import { fetchProject } from "@/queries/fetchProject.ts";

export const ProjectSelectionScreen = () => {
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
      <div className="w-full h-full flex flex-col items-center justify-center gap-6">
        <h1 className="text-3xl">Open project </h1>
        {availableProjects.length > 0 && (
          <div className="flex flex-col items-center justify-center gap-2">
            {availableProjects.map((p) => (
              <div
                onClick={handleOpenProject(p.id)}
                className="px-4 py-2 border gap-2 border-neutral-100 rounded-xl shadow hover:shadow-md transition-all cursor-pointer w-80 flex flex-col items-start gap-1"
                key={p.id}
              >
                <p className="font-semibold">{p.name}</p>
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
