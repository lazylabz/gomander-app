import { dataService } from "@/contracts/service.ts";
import type { Command } from "@/contracts/types.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { commandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

export const loadCommandDataIntoStore = (commands: Command[]) => {
  const { setCommands, setCommandsStatus, commandsStatus } =
    commandStore.getState();

  setCommands(commands);

  const newCommandsStatus: Record<string, CommandStatus> = {};
  commands.forEach((command) => {
    if (!newCommandsStatus[command.id]) {
      newCommandsStatus[command.id] =
        commandsStatus[command.id] || CommandStatus.IDLE;
    }
  });

  setCommandsStatus(newCommandsStatus);
};

export const fetchCommands = async () => {
  const commands = await dataService.getCommands();

  loadCommandDataIntoStore(commands);

  await fetchCommandGroups(); // Update command groups as they contain commands
};
