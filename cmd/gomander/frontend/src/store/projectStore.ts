import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

import type { Project } from "@/contracts/types.ts";

type ProjectStore = {
  projectInfo: Project | null;
  setProjectInfo: (project: Project | null) => void;
};

export const projectStore = createStore<ProjectStore>()((set) => ({
  projectInfo: null,
  setProjectInfo: (projectInfo: Project | null) => {
    set({ projectInfo: projectInfo });
  },
}));

export const useProjectStore = <T>(selector: (state: ProjectStore) => T): T => {
  return useStore(projectStore, selector);
};
