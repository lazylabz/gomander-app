import { dataService } from "@/contracts/service.ts";
import type { Command } from "@/contracts/types.ts";
import { commandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

export const loadCommandDataIntoStore = (
  commandsData: Record<string, Command>,
) => {
  const { setCommands, setCommandsStatus, commandsStatus } =
    commandStore.getState();

  setCommands(commandsData);

  const newCommandsStatus = Object.fromEntries(
    Object.keys(commandsData).map((id) => [
      id,
      commandsStatus[id] || CommandStatus.IDLE,
    ]),
  );

  setCommandsStatus(newCommandsStatus);
};

export const fetchCommands = async () => {
  const commandsData = await dataService.getCommands();

  loadCommandDataIntoStore(commandsData);
};
