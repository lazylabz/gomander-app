import { dataService } from "@/contracts/service.ts";
import { projectStore } from "@/store/projectStore.ts";

export const closeProject = async () => {
  const { setProject } = projectStore.getState();
  await dataService.closeProject();
  setProject(null);
};
