import { commandStore } from "@/store/commandStore";
import type { CommandStatus } from "@/types/CommandStatus.ts";

export const updateCommandStatus = (id: string, status: CommandStatus) => {
  const { commandsStatus, setCommandsStatus } = commandStore.getState();

  const updatedStatus = {
    ...commandsStatus,
    [id]: status,
  };

  setCommandsStatus(updatedStatus);
};
