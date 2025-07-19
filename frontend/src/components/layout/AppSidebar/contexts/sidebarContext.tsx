import { createContext, useContext } from "react";

import type { Command } from "@/types/contracts.ts";

export const sidebarContext = createContext<{
  editCommand: (command: Command) => void;
}>(
  {} as {
    editCommand: (command: Command) => void;
  },
);

export const useSidebarContext = () => {
  return useContext(sidebarContext);
};
