import { dataService } from "@/contracts/service.ts";

export const exportProject = async (projectId: string) => {
  await dataService.exportProject(projectId);
};
