import { dataService } from "@/contracts/service.ts";

export const exportProject = async (projectId: string) => {
  return dataService.exportProject(projectId);
};
