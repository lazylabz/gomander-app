import { dataService } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
import { projectStore } from "@/store/projectStore.ts";

export const loadAllProjectData = async () => {
  const { setCommands } = commandStore.getState();
  const { setCommandGroups } = commandGroupStore.getState();
  const { setProjectInfo } = projectStore.getState();

  const [commands, commandGroups, projectInfo] = await Promise.all([
    dataService.getCommands(),
    dataService.getCommandGroups(),
    dataService.getCurrentProject(),
  ]);

  console.log({
    commands,
    commandGroups,
    projectInfo,
  });

  setCommands(commands);
  setCommandGroups(commandGroups);
  setProjectInfo(projectInfo);
};
