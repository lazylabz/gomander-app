import { dataService } from "@/contracts/service.ts";
import type { CommandGroup } from "@/contracts/types.ts";

export const createCommandGroup = async (commandGroup: CommandGroup) => {
  await dataService.createCommandGroup(commandGroup);
};
