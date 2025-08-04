import { dataService } from "@/contracts/service.ts";

export const deleteProject = async (projectIdBeingDeleted: string) => {
  await dataService.deleteProject(projectIdBeingDeleted);
};