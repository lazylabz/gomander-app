import { dataService } from "@/contracts/service.ts";
import type { Project } from "@/contracts/types.ts";

export const editOpenedProject = async (project: Project) => {
  await dataService.editProject(project);
};
