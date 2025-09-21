import { zodResolver } from "@hookform/resolvers/zod";
import { createContext, useContext, useEffect, useState } from "react";
import { useForm, type UseFormReturn } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { useLocation, useNavigate } from "react-router";
import { toast } from "sonner";

import { useTheme } from "@/contexts/theme.tsx";
import { translationsService } from "@/contracts/service.ts";
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

type SupportedLanguage = {
  value: string;
  label: string;
};

const languageValueToLabelMap: Record<string, string> = {
  en: "English",
  es: "Español",
};

// Define context
export interface SettingsContextData {
  settingsForm: SettingsUseFormReturn;
  hasUnsavedChanges: boolean;
  initialTab: SettingsTab;
  supportedLanguages: SupportedLanguage[];
  saveSettings: (formData: SettingsFormType) => Promise<void>;
  closeSettings: () => void;
}

export const settingsContext = createContext<SettingsContextData>({
  settingsForm: {} as SettingsUseFormReturn,
  hasUnsavedChanges: false,
  initialTab: SettingsTab.User,
  supportedLanguages: [],
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
  const { i18n } = useTranslation();

  const [supportedLanguages, setSupportedLanguages] = useState<
    SupportedLanguage[]
  >([]);

  const navigate = useNavigate();
  const { state } = useLocation();

  const initialTab = state?.tab || SettingsTab.User;

  useEffect(() => {
    const loadSupportedLanguages = async () => {
      try {
        const languages = await translationsService.getSupportedLanguages();
        const languageOptions = languages.map(
          (lang): SupportedLanguage => ({
            value: lang,
            label: languageValueToLabelMap[lang] || lang,
          }),
        );
        setSupportedLanguages(languageOptions);
      } catch (error) {
        console.error("Failed to load supported languages:", error);
      }
    };

    loadSupportedLanguages();
  }, []);

  const settingsForm = useForm<SettingsFormType>({
    resolver: zodResolver(settingsFormSchema),
    values: {
      environmentPaths: userConfig.environmentPaths,
      theme: rawTheme,
      logLineLimit: userConfig.logLineLimit,
      locale: i18n.language,
      name: projectInfo?.name || "",
      baseWorkingDirectory: projectInfo?.workingDirectory || "",
    },
  });

  // Apply theme when it changes without needing to submit the form
  settingsForm.register("theme", {
    onChange: (e) => {
      const newTheme = e.target.value;
      setRawTheme(newTheme);
    },
  });

  const loadLanguage = async (lang: string) => {
    if (!i18n.hasResourceBundle(lang, "translation")) {
      const translations = await translationsService.getTranslation(lang);
      i18n.addResourceBundle(lang, "translation", translations);
    }

    await i18n.changeLanguage(lang);
  };

  const { dirtyFields } = settingsForm.formState;

  const hasUserChanges = !!dirtyFields.environmentPaths || !!dirtyFields.logLineLimit || !!dirtyFields.locale;
  const hasProjectChanges =
    dirtyFields.name || dirtyFields.baseWorkingDirectory;
  const saveSettings = async (formData: SettingsFormType) => {
    if (!projectInfo) {
      return;
    }

    // Save user settings
    if (hasUserChanges) {
      try {
        await loadLanguage(formData.locale);
        await saveUserConfig({
          lastOpenedProjectId: userConfig.lastOpenedProjectId,
          environmentPaths: formData.environmentPaths,
          logLineLimit: formData.logLineLimit,
          locale: formData.locale,
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
    supportedLanguages,
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
