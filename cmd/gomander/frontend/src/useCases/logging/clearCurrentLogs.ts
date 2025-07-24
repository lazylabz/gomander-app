import { commandStore } from "@/store/commandStore.ts";

export const clearCurrentLogs = () => {
  const { activeCommandId, commandsLogs, setCommandsLogs } =
    commandStore.getState();

  if (!activeCommandId) {
    return;
  }

  setCommandsLogs({
    ...commandsLogs,
    [activeCommandId]: [],
  });
};
