import { dataService } from "@/contracts/service.ts";
import type { Command } from "@/contracts/types.ts";

export const duplicateCommand = async (
  command: Command,
  targetGroupId?: string,
) => {
  await dataService.duplicateCommand(command.id, targetGroupId || "");
};
