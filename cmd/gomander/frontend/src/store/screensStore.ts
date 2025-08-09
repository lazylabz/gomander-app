import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

export enum Screen {
  Logs = "logs",
  ProjectSelection = "project-selection",
  Settings = "settings",
}

type ScreensStore = {
  currentScreen: Screen;
  setCurrentScreen: (screen: Screen) => void;
};

export const screensStore = createStore<ScreensStore>()((set) => ({
  currentScreen: Screen.Logs,
  setCurrentScreen: (screen: Screen) => set({ currentScreen: screen }),
}));

export const useScreensStore = <T>(selector: (state: ScreensStore) => T): T => {
  return useStore(screensStore, selector);
};
