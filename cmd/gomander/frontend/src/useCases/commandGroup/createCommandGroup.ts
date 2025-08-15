import { dataService } from "@/contracts/service.ts";
import type { CommandGroup } from "@/contracts/types.ts";
import { isDefined } from "@/helpers/mapHelpers.ts";
import { commandStore } from "@/store/commandStore.ts";

interface CreateCommandGroupParams extends Omit<CommandGroup, "commands"> {
  commands: string[];
}

export const createCommandGroup = async (args: CreateCommandGroupParams) => {
  const { commands } = commandStore.getState();

  const groupCommands = args.commands
    .map((commandId) => {
      const command = commands.find((c) => c.id === commandId);
      if (!command) {
        return undefined;
      }
      return command;
    })
    .filter(isDefined);

  const commandGroup: CommandGroup = {
    id: args.id,
    projectId: args.projectId,
    name: args.name,
    commands: groupCommands,
    position: 0, // Will be set by the backend
  };
  await dataService.createCommandGroup(commandGroup);
};
