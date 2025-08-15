import { toast } from "sonner";

import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { startCommand } from "@/useCases/command/startCommand.ts";

export const runCommandGroup = async (groupId: string) => {
  const { commandGroups } = commandGroupStore.getState();
  const { commandsStatus } = commandStore.getState();

  const group = commandGroups.find((g) => g.id === groupId);
  if (!group) {
    toast.error("Command group not found");
    return;
  }

  const notRunningCommands = group.commands.filter(
    (cmd) => commandsStatus[cmd.id] !== CommandStatus.RUNNING,
  );

  await Promise.all(
    notRunningCommands.map(async (cmd) => {
      await startCommand(cmd.id);
    }),
  );
};
