import { dataService } from "@/contracts/service.ts";
import type { Command } from "@/contracts/types.ts";

export const createCommand = async (command: Command) => {
  await dataService.addCommand(command);
};
