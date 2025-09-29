import { createContext, useContext } from "react";

import type { Command } from "@/contracts/types.ts";

export const sidebarContext = createContext<{
  startEditingCommand: (command: Command) => void;
}>(
  {} as {
    startEditingCommand: (command: Command) => void;
  },
);

export const useSidebarContext = () => {
  return useContext(sidebarContext);
};
