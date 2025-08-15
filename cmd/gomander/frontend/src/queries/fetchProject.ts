import { dataService } from "@/contracts/service.ts";
import { projectStore } from "@/store/projectStore.ts";

export const fetchProject = async () => {
  const { setProjectInfo } = projectStore.getState();

  const projectInfo = await dataService.getCurrentProject();

  setProjectInfo(projectInfo);
};
