import { commandStore } from "@/store/commandStore.ts";
import { cleanCommandLogs } from "@/useCases/command/cleanCommandLogs.ts";

export const clearCurrentLogs = () => {
  const { activeCommandId } = commandStore.getState();

  if (!activeCommandId) {
    return;
  }

  cleanCommandLogs(activeCommandId);
};
