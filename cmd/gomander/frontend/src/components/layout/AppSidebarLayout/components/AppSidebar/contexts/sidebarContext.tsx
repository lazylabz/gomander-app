import { createContext, useContext } from "react";

import type { Command, CommandGroup } from "@/contracts/types.ts";

export const sidebarContext = createContext<{
  startEditingCommand: (command: Command) => void;
  startEditingCommandGroup: (commandGroup: CommandGroup) => void;
  isReorderingGroups: boolean;
}>(
  {} as {
    startEditingCommand: (command: Command) => void;
    startEditingCommandGroup: (commandGroup: CommandGroup) => void;
    isReorderingGroups: boolean;
  },
);

export const useSidebarContext = () => {
  return useContext(sidebarContext);
};
