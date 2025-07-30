import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

import type { Project } from "@/contracts/types.ts";

type ProjectStore = {
  project: Project | null;
  setProject: (project: Project | null) => void;
};

export const projectStore = createStore<ProjectStore>()((set) => ({
  project: null,
  setProject: (project: Project | null) => {
    set({ project });
  },
}));

export const useProjectStore = <T>(selector: (state: ProjectStore) => T): T => {
  return useStore(projectStore, selector);
};
