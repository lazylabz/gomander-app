import { zodResolver } from "@hookform/resolvers/zod";
import { createContext, useContext } from "react";
import { useForm, type UseFormReturn } from "react-hook-form";
import { useLocation, useNavigate } from "react-router";

import { type Theme, useTheme } from "@/contexts/theme.tsx";
import { fetchProject } from "@/queries/fetchProject.ts";
import { useProjectStore } from "@/store/projectStore.ts";
import { useUserConfigurationStore } from "@/store/userConfigurationStore.ts";
import { editOpenedProject } from "@/useCases/project/editOpenedProject.ts";
import { saveUserConfig } from "@/useCases/userConfig/saveUserConfig.ts";

import {
  settingsFormSchema,
  type SettingsFormType,
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
  unknown,
  SettingsFormData
>;

// Define context
export interface SettingsContextData {
  initialTab: SettingsTab;
  closeSettings: () => void;
  saveSettings: (formData: SettingsFormType) => Promise<void>;
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
  const { rawTheme, setRawTheme } = useTheme();

  const navigate = useNavigate();
  const { state } = useLocation();

  const initialTab = state?.tab || SettingsTab.User;

  const settingsForm = useForm<SettingsFormType>({
    resolver: zodResolver(settingsFormSchema),
    values: {
      environmentPaths: userConfig.environmentPaths.map((p) => ({ value: p })),
      theme: rawTheme,
      name: projectInfo?.name || "",
      baseWorkingDirectory: projectInfo?.baseWorkingDirectory || "",
    },
  });

  const closeSettings = () => {
    navigate(-1);
  };

  const saveSettings = async (formData: SettingsFormType) => {
    // Save user settings
    setRawTheme(formData.theme);
    await saveUserConfig({
      lastOpenedProjectId: userConfig.lastOpenedProjectId,
      environmentPaths: formData.environmentPaths.map((path) => path.value),
    });

    // Save project settings
    if (!projectInfo) {
      return;
    }
    await editOpenedProject({
      ...projectInfo,
      name: formData.name,
      baseWorkingDirectory: formData.baseWorkingDirectory,
    });
    await fetchProject();
  };

  const value: SettingsContextData = {
    initialTab,
    closeSettings,
    saveSettings,
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
