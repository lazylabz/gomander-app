import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

import type { UserConfig } from "@/contracts/types.ts";

type UserConfigurationStore = {
  userConfig: UserConfig;
  setUserConfig: (config: UserConfig) => void;
  isLoaded: boolean;
};

export const userConfigurationStore = createStore<UserConfigurationStore>()(
  (set) => ({
    isLoaded: false,
    userConfig: {
      environmentPaths: [],
      lastOpenedProjectId: "",
      logLineLimit: 100,
      locale: "en",
    },
    setUserConfig: (config: UserConfig) => {
      set({ userConfig: config, isLoaded: true });
    },
  }),
);

export const useUserConfigurationStore = <T>(
  selector: (state: UserConfigurationStore) => T,
): T => {
  return useStore(userConfigurationStore, selector);
};
