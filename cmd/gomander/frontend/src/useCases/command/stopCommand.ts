import { dataService } from "@/contracts/service.ts";

export const stopCommand = async (commandId: string) => {
  await dataService.stopCommand(commandId);
};
