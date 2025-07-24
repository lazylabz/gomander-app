import { dataService } from "@/contracts/service.ts";
import { commandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

export const startCommand = async (commandId: string) => {
  const { setCommandsStatus, commandsStatus, commandsLogs, setCommandsLogs } =
    commandStore.getState();

  const newCommandsStatus = {
    ...commandsStatus,
    [commandId]: CommandStatus.RUNNING,
  };

  const newCommandsLogs = {
    ...commandsLogs,
    [commandId]: [], // Initialize logs for the command being executed
  };

  setCommandsStatus(newCommandsStatus);
  setCommandsLogs(newCommandsLogs);

  await dataService.runCommand(commandId);
};
