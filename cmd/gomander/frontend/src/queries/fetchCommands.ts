import { dataService } from "@/contracts/service.ts";
import { commandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

export const fetchCommands = async () => {
  const { setCommands, setCommandsStatus, commandsStatus } =
    commandStore.getState();

  const commandsData = await dataService.getCommands();

  setCommands(commandsData);

  const newCommandsStatus = Object.fromEntries(
    Object.keys(commandsData).map((id) => [
      id,
      commandsStatus[id] || CommandStatus.IDLE,
    ]),
  );

  setCommandsStatus(newCommandsStatus);
};
