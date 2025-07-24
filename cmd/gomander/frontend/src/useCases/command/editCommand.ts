import { dataService } from "@/contracts/service.ts";
import type { Command } from "@/contracts/types.ts";

export const editCommand = async (command: Command) => {
  await dataService.editCommand(command);
};
