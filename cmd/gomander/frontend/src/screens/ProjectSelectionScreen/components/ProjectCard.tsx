import { EllipsisVertical } from "lucide-react";
import { toast } from "sonner";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { dataService } from "@/contracts/service.ts";
import type { Project } from "@/contracts/types.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { loadAllProjectData } from "@/queries/loadAllProjectData.ts";
import { exportProject } from "@/useCases/project/exportProject.ts";

export const ProjectCard = ({
  project,
  handleDeleteProject,
}: {
  project: Project;
  handleDeleteProject: () => void;
}) => {
  const handleOpenProject = (projectId: string) => async () => {
    await dataService.openProject(projectId);
    await loadAllProjectData();
  };

  const handleExportProject = async () => {
    try {
      await exportProject(project.id);
      toast.success("Project exported successfully!");
    } catch (e) {
      toast.error(parseError(e, "Failed to export the project"));
    }
  };

  return (
    <div className="relative px-4 py-2 border border-neutral-100 dark:border-neutral-900 rounded-xl shadow-none hover:shadow-md shadow-neutral-100 dark:shadow-neutral-800 transition-all w-80">
      <DropdownMenu>
        <DropdownMenuTrigger className="cursor-pointer flex absolute items-center justify-center top-0 right-0 px-2 pb-2 pt-3 text-muted-foreground hover:text-primary">
          <EllipsisVertical size={16} />
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem
            className="cursor-pointer"
            onClick={handleDeleteProject}
          >
            Delete
          </DropdownMenuItem>
          <DropdownMenuItem
            className="cursor-pointer"
            onClick={handleExportProject}
          >
            Export
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <div
        role="button"
        className="flex flex-col items-start gap-1 p-2 cursor-pointer"
        onClick={handleOpenProject(project.id)}
      >
        <p className="font-medium">{project.name}</p>
      </div>
    </div>
  );
};
