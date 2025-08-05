import { dataService } from "@/contracts/service.ts";
import type { Project } from "@/contracts/types.ts";

export const importProject = async (project: Project) => {
  await dataService.importProject(project);
};
