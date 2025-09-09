import { dataService } from "@/contracts/service.ts";

export const stopCommandGroup = async (groupId: string) => {
  await dataService.stopCommandGroup(groupId);
};
