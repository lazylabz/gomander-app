import { createContext, useContext } from "react";

import type { Command, CommandGroup } from "@/contracts/types.ts";

export const sidebarContext = createContext<{
  editCommand: (command: Command) => void;
  editCommandGroup: (commandGroup: CommandGroup) => void;
}>(
  {} as {
    editCommand: (command: Command) => void;
    editCommandGroup: (commandGroup: CommandGroup) => void;
  },
);

export const useSidebarContext = () => {
  return useContext(sidebarContext);
};
