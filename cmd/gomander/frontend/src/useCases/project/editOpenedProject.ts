import { dataService } from "@/contracts/service.ts";
import type { ProjectInfo } from "@/contracts/types.ts";

export const editOpenedProject = async (projectInfo: ProjectInfo) => {
  await dataService.editProject(projectInfo);
};
