import { dataService } from "@/contracts/service.ts";
import { cleanCommandError } from "@/useCases/command/cleanCommandError.ts";

export const startCommand = async (commandId: string) => {
  cleanCommandError(commandId);

  await dataService.runCommand(commandId);
};
