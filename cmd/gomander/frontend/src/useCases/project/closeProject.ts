import { dataService } from "@/contracts/service.ts";
import { projectStore } from "@/store/projectStore.ts";

export const closeProject = async () => {
  const { setProjectInfo } = projectStore.getState();
  await dataService.closeProject();
  setProjectInfo(null);
};
