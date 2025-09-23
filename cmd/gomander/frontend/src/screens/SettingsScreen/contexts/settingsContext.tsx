import { zodResolver } from "@hookform/resolvers/zod";
import { createContext, useContext } from "react";
import { useForm, type UseFormReturn } from "react-hook-form";
import { useLocation, useNavigate } from "react-router";
import { toast } from "sonner";

import { useTheme } from "@/contexts/theme.tsx";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchProject } from "@/queries/fetchProject.ts";
import { fetchUserConfig } from "@/queries/fetchUserConfig.ts";
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

type SettingsUseFormReturn = UseFormReturn<
  SettingsFormType,
  unknown,
  SettingsFormType
>;

// Define context
export interface SettingsContextData {
  settingsForm: SettingsUseFormReturn;
  hasUnsavedChanges: boolean;
  initialTab: SettingsTab;
  saveSettings: (formData: SettingsFormType) => Promise<void>;
  closeSettings: () => void;
}

export const settingsContext = createContext<SettingsContextData>({
  settingsForm: {} as SettingsUseFormReturn,
  hasUnsavedChanges: false,
  initialTab: SettingsTab.User,
  saveSettings: async () => {},
  closeSettings: () => {},
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

  // Initialize form including failurePatterns
  const settingsForm = useForm<SettingsFormType>({
    resolver: zodResolver(settingsFormSchema),
    values: {
      environmentPaths: userConfig.environmentPaths,
      theme: rawTheme,
      logLineLimit: userConfig.logLineLimit,
      name: projectInfo?.name || "",
      baseWorkingDirectory: projectInfo?.workingDirectory || "",
      failurePatterns: projectInfo?.failurePatterns || [], 
    },
  });

  // Apply theme when it changes without needing to submit the form
  settingsForm.register("theme", {
    onChange: (e) => {
      const newTheme = e.target.value;
      setRawTheme(newTheme);
    },
  });

  const { dirtyFields } = settingsForm.formState;

  const hasUserChanges = !!dirtyFields.environmentPaths || !!dirtyFields.logLineLimit;
  const hasProjectChanges =
    dirtyFields.name || dirtyFields.baseWorkingDirectory || dirtyFields.failurePatterns; // UPDATED

  const saveSettings = async (formData: SettingsFormType) => {
    if (!projectInfo) return;

    // Save user settings
    if (hasUserChanges) {
      try {
        await saveUserConfig({
          lastOpenedProjectId: userConfig.lastOpenedProjectId,
          environmentPaths: formData.environmentPaths,
          logLineLimit: formData.logLineLimit,
        });
        toast.success("User settings saved successfully");
      } catch (e) {
        throw new Error(parseError(e, "Failed to save user settings"));
      }
      await fetchUserConfig();
    }

    // Save project settings
    if (hasProjectChanges) {
      try {
        await editOpenedProject({
          ...projectInfo,
          name: formData.name,
          workingDirectory: formData.baseWorkingDirectory,
          failurePatterns: formData.failurePatterns, 
        });
        toast.success("Project settings saved successfully");
      } catch (e) {
        throw new Error(parseError(e, "Failed to save project settings"));
      }

      await fetchProject();
    }

    settingsForm.reset();
  };

  const closeSettings = () => {
    settingsForm.reset();
    navigate(-1);
  };

  const value: SettingsContextData = {
    settingsForm,
    hasUnsavedChanges: settingsForm.formState.isDirty,
    initialTab,
    saveSettings,
    closeSettings,
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
