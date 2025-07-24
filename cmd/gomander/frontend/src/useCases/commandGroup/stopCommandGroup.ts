import { toast } from "sonner";

import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus";
import { stopCommand } from "@/useCases/command/stopCommand.ts";

export const stopCommandGroup = async (groupId: string) => {
  const { commandGroups } = commandGroupStore.getState();
  const { commandsStatus } = commandStore.getState();

  const group = commandGroups.find((g) => g.id === groupId);
  if (!group) {
    toast.error("Command group not found");
    return;
  }

  const runningCommands = group.commands.filter(
    (cmdId) => commandsStatus[cmdId] === CommandStatus.RUNNING,
  );

  console.log(groupId);

  await Promise.all(
    runningCommands.map(async (cmdId) => {
      await stopCommand(cmdId);
    }),
  );
};
