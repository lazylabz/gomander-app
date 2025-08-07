import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

import type { ProjectInfo } from "@/contracts/types.ts";

type ProjectStore = {
  projectInfo: ProjectInfo | null;
  setProjectInfo: (project: ProjectInfo | null) => void;
};

export const projectStore = createStore<ProjectStore>()((set) => ({
  projectInfo: null,
  setProjectInfo: (projectInfo: ProjectInfo | null) => {
    set({ projectInfo });
  },
}));

export const useProjectStore = <T>(selector: (state: ProjectStore) => T): T => {
  return useStore(projectStore, selector);
};
