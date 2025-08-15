import { dataService } from "@/contracts/service.ts";
import type { CommandGroup } from "@/contracts/types.ts";

export const reorderCommandGroupCommands = async (
  commandGroup: CommandGroup,
) => {
  await dataService.editCommandGroup(commandGroup);
};
