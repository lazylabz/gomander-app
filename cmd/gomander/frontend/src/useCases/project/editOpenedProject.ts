import { dataService } from "@/contracts/service.ts";
import type { Project } from "@/contracts/types.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";

type ProjectBasicData = Pick<Project, "id" | "name" | "baseWorkingDirectory">;

export const editOpenedProject = async (project: ProjectBasicData) => {
  const { commands } = commandStore.getState();
  const { commandGroups } = commandGroupStore.getState();

  await dataService.editProject({
    ...project,
    commands,
    commandGroups,
  });
};
