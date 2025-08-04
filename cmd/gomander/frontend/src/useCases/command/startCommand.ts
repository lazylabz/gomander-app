import { dataService } from "@/contracts/service.ts";
import { cleanCommandLogs } from "@/useCases/command/cleanCommandLogs.ts";

export const startCommand = async (commandId: string) => {
  cleanCommandLogs(commandId);

  await dataService.runCommand(commandId);
};
