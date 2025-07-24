import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

import type { CommandGroup } from "@/contracts/types.ts";

type CommandGroupStore = {
  commandGroups: CommandGroup[];
  setCommandGroups: (groups: CommandGroup[]) => void;
};

export const commandGroupStore = createStore<CommandGroupStore>()((set) => ({
  commandGroups: [],
  setCommandGroups: (groups) => set({ commandGroups: groups }),
}));

export const useCommandGroupStore = <T>(
  selector: (state: CommandGroupStore) => T,
): T => {
  return useStore(commandGroupStore, selector);
};
