import { dataService } from "@/contracts/service.ts";

export const startCommand = async (commandId: string) => {
  await dataService.runCommand(commandId);
};
