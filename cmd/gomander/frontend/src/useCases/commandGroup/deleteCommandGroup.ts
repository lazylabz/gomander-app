import { dataService } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";

export const deleteCommandGroup = async (commandGroupId: string) => {
  const { commandGroups } = commandGroupStore.getState();

  const updatedCommandGroups = commandGroups.filter(
    (cg) => cg.id !== commandGroupId,
  );

  await dataService.saveCommandGroups(updatedCommandGroups);
};
