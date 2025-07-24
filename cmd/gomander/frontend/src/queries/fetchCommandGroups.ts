import { dataService } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";

export const fetchCommandGroups = async (): Promise<void> => {
  const { setCommandGroups } = commandGroupStore.getState();
  const groups = await dataService.getCommandGroups();
  setCommandGroups(groups);
};
