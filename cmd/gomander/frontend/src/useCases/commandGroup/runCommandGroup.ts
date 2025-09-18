import { dataService } from "@/contracts/service.ts";

export const runCommandGroup = async (groupId: string) => {
  await dataService.runCommandGroup(groupId);
};
