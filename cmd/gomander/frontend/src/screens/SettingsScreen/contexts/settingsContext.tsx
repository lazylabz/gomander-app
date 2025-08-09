import { zodResolver } from "@hookform/resolvers/zod";
import { createContext, useContext } from "react";
import { useForm, type UseFormReturn } from "react-hook-form";
import { useLocation, useNavigate } from "react-router";

import { type Theme, useTheme } from "@/contexts/theme.tsx";
import { useProjectStore } from "@/store/projectStore.ts";
import { useUserConfigurationStore } from "@/store/userConfigurationStore.ts";

import {
  settingsFormSchema,
  type SettingsFormSchemaType,
} from "./settingsFormSchema";

export enum SettingsTab {
  User = "user",
  Project = "project",
}

interface SettingsFormData {
  // User settings
  environmentPaths: { value: string }[];
  theme: Theme;
  // Project settings
  name: string;
  baseWorkingDirectory: string;
}

type SettingsUseFormReturn = UseFormReturn<
  SettingsFormData,
  never,
  SettingsFormData
>;

// Define context
export interface SettingsContextData {
  initialTab: SettingsTab;
  closeSettings: () => void;
  saveSettings: (closeOnSave: boolean) => Promise<void>;
  hasUnsavedChanges: boolean;
  settingsForm: SettingsUseFormReturn;
}

export const settingsContext = createContext<SettingsContextData>({
  initialTab: SettingsTab.User,
  closeSettings: () => {},
  saveSettings: async () => {},
  hasUnsavedChanges: false,
  settingsForm: {} as SettingsUseFormReturn,
});

// Define provider
export const SettingsContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const userConfig = useUserConfigurationStore((state) => state.userConfig);
  const projectInfo = useProjectStore((state) => state.projectInfo);
  const { rawTheme } = useTheme();

  const navigate = useNavigate();
  const { state } = useLocation();

  const initialTab = state?.tab || SettingsTab.User;

  const closeSettings = () => {
    navigate(-1);
  };

  const settingsForm = useForm<SettingsFormSchemaType>({
    resolver: zodResolver(settingsFormSchema),
    values: {
      environmentPaths: userConfig.environmentPaths.map((p) => ({ value: p })),
      theme: rawTheme,
      name: projectInfo?.name || "",
      baseWorkingDirectory: projectInfo?.baseWorkingDirectory || "",
    },
  });

  const value: SettingsContextData = {
    initialTab,
    closeSettings,
    saveSettings: async (closeOnSave) => {
      // TODO: Logic to save settings
      if (closeOnSave) {
        // TODO: Logic to close settings after saving
      }
    },
    hasUnsavedChanges: false, // TODO: Implement logic to track unsaved changes
    settingsForm,
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
