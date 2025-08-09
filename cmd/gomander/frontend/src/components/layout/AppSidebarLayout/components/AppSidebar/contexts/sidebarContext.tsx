import { createContext, useContext } from "react";

import type { Command, CommandGroup } from "@/contracts/types.ts";

export const sidebarContext = createContext<{
  startEditingCommand: (command: Command) => void;
  startEditingCommandGroup: (commandGroup: CommandGroup) => void;
}>(
  {} as {
    startEditingCommand: (command: Command) => void;
    startEditingCommandGroup: (commandGroup: CommandGroup) => void;
  },
);

export const useSidebarContext = () => {
  return useContext(sidebarContext);
};
