import { dataService } from "@/contracts/service.ts";
import type { CommandGroup } from "@/contracts/types.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";

export const editCommandGroup = async (commandGroup: CommandGroup) => {
  const { commandGroups } = commandGroupStore.getState();

  const updatedCommandGroups = commandGroups.map((cg) =>
    cg.id === commandGroup.id ? commandGroup : cg,
  );

  await dataService.saveCommandGroups(updatedCommandGroups);
};
