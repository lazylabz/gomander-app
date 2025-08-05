import { EllipsisVertical } from "lucide-react";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { dataService } from "@/contracts/service.ts";
import type { Project } from "@/contracts/types.ts";
import { fetchProject } from "@/queries/fetchProject.ts";
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
    await fetchProject();
  };

  return (
    <div className="relative px-4 py-2 border border-neutral-100 rounded-xl shadow hover:shadow-md transition-all w-80">
      <DropdownMenu>
        <DropdownMenuTrigger className="absolute top-2 right-2 ">
          <EllipsisVertical
            size={16}
            className="text-muted-foreground cursor-pointer hover:text-primary"
          />
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem onClick={handleDeleteProject}>
            Delete
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => exportProject(project.id)}>
            Export
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <div
        role="button"
        className="flex flex-col items-start gap-1 cursor-pointer"
        onClick={handleOpenProject(project.id)}
      >
        <p className="font-semibold">{project.name}</p>
        <div className="flex text-sm gap-4 text-muted-foreground">
          <p>{Object.keys(project.commands).length} commands</p>
          <p>{project.commandGroups.length} command groups</p>
        </div>
      </div>
    </div>
  );
};
