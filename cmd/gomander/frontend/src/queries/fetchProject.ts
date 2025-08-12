import { dataService } from "@/contracts/service.ts";
import type { Project } from "@/contracts/types.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { projectStore } from "@/store/projectStore.ts";

const loadProjectDataIntoStores = (project: Project) => {
  const { setCommandGroups } = commandGroupStore.getState();
  const { setProjectInfo } = projectStore.getState();

  setProjectInfo({
    id: project.id,
    name: project.name,
    baseWorkingDirectory: project.baseWorkingDirectory,
  });
  setCommandGroups(project.commandGroups);
};
export const fetchProject = async () => {
  const project = await dataService.getCurrentProject();

  console.log({ project });

  if (project) {
    await fetchCommands();
    loadProjectDataIntoStores(project);
  }
};
