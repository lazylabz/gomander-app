import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

import type { Command } from "@/contracts/types.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

type CommandStore = {
  commands: Record<string, Command>;
  setCommands: (commands: Record<string, Command>) => void;

  commandsStatus: Record<string, CommandStatus>;
  setCommandsStatus: (commandsStatus: Record<string, CommandStatus>) => void;

  activeCommandId: string | null;
  setActiveCommandId: (commandId: string | null) => void;

  commandsLogs: Record<string, string[]>;
  setCommandsLogs: (logs: Record<string, string[]>) => void;
};

// To be used in use cases
export const commandStore = createStore<CommandStore>()((set) => ({
  commands: {},
  setCommands: (commands) => set({ commands }),

  commandsStatus: {},
  setCommandsStatus: (commandsStatus) => set({ commandsStatus }),

  activeCommandId: null,
  setActiveCommandId: (commandId) => set({ activeCommandId: commandId }),

  commandsLogs: {},
  setCommandsLogs: (logs) => set({ commandsLogs: logs }),
}));

// To be used in react components
export const useCommandStore = <T>(selector: (state: CommandStore) => T): T =>
  useStore(commandStore, selector);
