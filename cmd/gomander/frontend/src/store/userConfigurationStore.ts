import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

import type { UserConfig } from "@/contracts/types.ts";

type UserConfigurationStore = {
  userConfig: UserConfig;
  setUserConfig: (config: UserConfig) => void;
};

export const userConfigurationStore = createStore<UserConfigurationStore>()(
  (set) => ({
    userConfig: {
      extraPaths: [],
      lastOpenedProjectId: "",
    },
    setUserConfig: (config: UserConfig) => {
      set({ userConfig: config });
    },
  }),
);

export const useUserConfigurationStore = <T>(
  selector: (state: UserConfigurationStore) => T,
): T => {
  return useStore(userConfigurationStore, selector);
};
