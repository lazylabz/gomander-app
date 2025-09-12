import { dataService } from "@/contracts/service.ts";

export const deleteCommandGroup = async (commandGroupId: string) => {
  await dataService.deleteCommandGroup(commandGroupId);
};
