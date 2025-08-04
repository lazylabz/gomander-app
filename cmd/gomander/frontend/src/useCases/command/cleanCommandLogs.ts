import { commandStore } from "@/store/commandStore.ts";

export const cleanCommandLogs = (commandId: string): void => {
  const { commandsLogs, setCommandsLogs } = commandStore.getState();

  setCommandsLogs({
    ...commandsLogs,
    [commandId]: [],
  });
};
