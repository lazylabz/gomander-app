import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

import type { Command } from "@/contracts/types.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";

type CommandStore = {
  commands: Command[];
  setCommands: (commands: Command[]) => void;

  commandsStatus: Record<string, CommandStatus>;
  setCommandsStatus: (commandsStatus: Record<string, CommandStatus>) => void;

  commandIdsWithErrors: string[];
  setCommandIdsWithErrors: (commandsWithErrors: string[]) => void;

  activeCommandId: string | null;
  setActiveCommandId: (commandId: string | null) => void;

  commandsLogs: Record<string, string[]>;
  setCommandsLogs: (logs: Record<string, string[]>) => void;
  addLogs: (logs: Map<string, string[]>, linesLimit?: number) => void;
};

// To be used in use cases
export const commandStore = createStore<CommandStore>()((set) => ({
  commands: [],
  setCommands: (commands) => set({ commands }),

  commandsStatus: {},
  setCommandsStatus: (commandsStatus) => set({ commandsStatus }),

  commandIdsWithErrors: [],
  setCommandIdsWithErrors: (commandsWithErrors) =>
    set({ commandIdsWithErrors: commandsWithErrors }),

  activeCommandId: null,
  setActiveCommandId: (commandId) => set({ activeCommandId: commandId }),

  commandsLogs: {},
  setCommandsLogs: (logs) => set({ commandsLogs: logs }),

  addLogs: (logs: Map<string, string[]>, linesLimit: number = 100) =>
    set((state) => {
      const newCommandsLogs = { ...state.commandsLogs };
      logs.forEach((lines, commandId) => {
        if (!newCommandsLogs[commandId]) {
          newCommandsLogs[commandId] = [];
        }
        newCommandsLogs[commandId] = [
          ...newCommandsLogs[commandId],
          ...lines,
        ].slice(-linesLimit); // Keep only the last `linesLimit` lines
      });
      return {
        ...state,
        commandsLogs: newCommandsLogs,
      };
    }),
}));

// To be used in react components
export const useCommandStore = <T>(selector: (state: CommandStore) => T): T =>
  useStore(commandStore, selector);
