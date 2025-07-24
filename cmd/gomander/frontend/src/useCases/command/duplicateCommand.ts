import { dataService } from "@/contracts/service.ts";
import type { Command } from "@/contracts/types.ts";

export const duplicateCommand = async (command: Command) => {
  const newCommand = {
    ...command,
    name: `${command.name} (copy)`,
    id: crypto.randomUUID(),
  };
  await dataService.addCommand(newCommand);
};
