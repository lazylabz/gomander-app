import { dataService } from "@/contracts/service.ts";

export const deleteCommand = async (commandId: string) => {
  await dataService.removeCommand(commandId);
};
