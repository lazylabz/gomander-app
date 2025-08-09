import { createContext, useContext } from "react";
import type { UseFormReturn } from "react-hook-form";

export enum SettingsTab {
  User = "user",
  Project = "project",
}

interface SettingsFormData {
  environmentPaths: { value: string }[];
  theme: string;
  name: string; // Project name
  baseWorkingDirectory: string;
}

// Define context
export interface SettingsContextData {
  initialTab: SettingsTab;
  closeSettings: () => void;
  saveSettings: (closeOnSave: boolean) => Promise<void>;
  hasUnsavedChanges: boolean;
  settingsForm: UseFormReturn<SettingsFormData>;
}

export const settingsContext = createContext<SettingsContextData>({
  initialTab: SettingsTab.User,
  closeSettings: () => {},
  saveSettings: async () => {},
  hasUnsavedChanges: false,
  settingsForm: {} as UseFormReturn<SettingsFormData>,
});

// Define provider
export const SettingsContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const value: SettingsContextData = {
    initialTab: SettingsTab.User, // TODO: Set this based on state.tab or default value
    closeSettings: () => {
      // TODO: Logic to close settings
    },
    saveSettings: async (closeOnSave) => {
      // TODO: Logic to save settings
      if (closeOnSave) {
        // TODO: Logic to close settings after saving
      }
    },
    hasUnsavedChanges: false, // TODO: Implement logic to track unsaved changes
    settingsForm: {} as UseFormReturn<SettingsFormData>, // TODO
  };

  return (
    <settingsContext.Provider value={value}>
      {children}
    </settingsContext.Provider>
  );
};

// Custom hook to use the settings context
export const useSettingsContext = () => {
  return useContext(settingsContext);
};
