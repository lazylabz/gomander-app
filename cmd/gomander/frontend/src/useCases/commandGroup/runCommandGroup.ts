import { dataService } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { cleanCommandLogs } from "@/useCases/command/cleanCommandLogs.ts";

export const runCommandGroup = async (groupId: string) => {
  const { commandGroups } = commandGroupStore.getState();
  commandGroups
    .find((g) => g.id === groupId)
    ?.commands.forEach((c) => cleanCommandLogs(c.id));

  await dataService.runCommandGroup(groupId);
};
