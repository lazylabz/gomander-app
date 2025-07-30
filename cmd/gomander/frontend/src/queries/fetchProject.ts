import { dataService } from "@/contracts/service.ts";
import type { Project } from "@/contracts/types.ts";
import { loadCommandDataIntoStore } from "@/queries/fetchCommands.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { projectStore } from "@/store/projectStore.ts";

const loadProjectDataIntoStores = (project: Project) => {
  const { setCommandGroups } = commandGroupStore.getState();
  const { setProject } = projectStore.getState();
  
  setProject(project);
  setCommandGroups(project.commandGroups);
  loadCommandDataIntoStore(project.commands);
};
export const fetchProject = async () => {
  const project = await dataService.getCurrentProject();

  console.log({ project });

  if (project) {
    loadProjectDataIntoStores(project);
  }
};
