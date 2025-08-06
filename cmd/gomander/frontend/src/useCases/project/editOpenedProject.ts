import { dataService } from "@/contracts/service.ts";
import type { ProjectInfo } from "@/contracts/types.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";

export const editOpenedProject = async (project: ProjectInfo) => {
  const { commands } = commandStore.getState();
  const { commandGroups } = commandGroupStore.getState();

  await dataService.editProject({
    ...project,
    commands,
    commandGroups,
  });
};
