import { zodResolver } from "@hookform/resolvers/zod";
import { createContext, useContext, useEffect, useRef, useState } from "react";
import { useForm, type UseFormReturn } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { useLocation } from "react-router";

import { translationsService } from "@/contracts/service.ts";
import {
  projectSettingsSchema,
  type ProjectSettingsSchemaType,
} from "@/screens/SettingsScreen/schemas/projectSettingsSchema.ts";
import {
  userSettingsSchema,
  type UserSettingsSchemaType,
} from "@/screens/SettingsScreen/schemas/userSettingsSchema.ts";
import { saveProjectSettingsForm } from "@/screens/SettingsScreen/useCases/saveProjectSettingsForm.ts";
import { saveUserSettingsForm } from "@/screens/SettingsScreen/useCases/saveUserSettingsForm.ts";
import { useProjectStore } from "@/store/projectStore.ts";
import { useUserConfigurationStore } from "@/store/userConfigurationStore.ts";

export enum SettingsTab {
  User = "user",
  Project = "project",
}

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
  initialTab: SettingsTab;
  hasPendingChanges: boolean;
  projectSettingsForm: UseFormReturn<ProjectSettingsSchemaType>;
  userSettingsForm: UseFormReturn<UserSettingsSchemaType>;
  supportedLanguages: SupportedLanguage[];
}

export const settingsContext = createContext<SettingsContextData>({
  initialTab: SettingsTab.User,
  hasPendingChanges: false,
  projectSettingsForm: {} as UseFormReturn<ProjectSettingsSchemaType>,
  userSettingsForm: {} as UseFormReturn<UserSettingsSchemaType>,
  supportedLanguages: [],
});

// Define provider
export const SettingsContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const { state } = useLocation();

  const initialTab = state?.tab || SettingsTab.User;

  //region Project settings

  const projectInfo = useProjectStore((state) => state.projectInfo);

  const [hasProjectPendingChanges, setHasProjectPendingChanges] =
    useState(false);
  const lastProjectSavedValues = useRef<ProjectSettingsSchemaType | null>(null);

  const projectForm = useForm<ProjectSettingsSchemaType>({
    resolver: zodResolver(projectSettingsSchema),
    defaultValues: {
      name: projectInfo?.name || "",
      baseWorkingDirectory: projectInfo?.workingDirectory || "",
    },
  });

  const projectFormWatcher = projectForm.watch();

  // Autosave project settings
  useEffect(() => {
    const currentValues = projectForm.getValues();

    if (lastProjectSavedValues.current === null) {
      lastProjectSavedValues.current = JSON.parse(
        JSON.stringify(currentValues),
      );
      return;
    }

    if (
      JSON.stringify(currentValues) ===
      JSON.stringify(lastProjectSavedValues.current)
    ) {
      return;
    }

    setHasProjectPendingChanges(true);
    const timeout = setTimeout(async () => {
      lastProjectSavedValues.current = JSON.parse(
        JSON.stringify(currentValues),
      );
      await projectForm.handleSubmit(saveProjectSettingsForm)();
      setHasProjectPendingChanges(false);
    }, 300);

    return () => clearTimeout(timeout);
  }, [projectForm, projectFormWatcher]);

  //endregion

  //region User settings

  const userConfig = useUserConfigurationStore((state) => state.userConfig);
  const { i18n } = useTranslation();

  const [hasUserPendingChanges, setHasUserPendingChanges] = useState(false);
  const lastUserSavedValues = useRef<UserSettingsSchemaType | null>(null);

  const userForm = useForm<UserSettingsSchemaType>({
    resolver: zodResolver(userSettingsSchema),
    defaultValues: {
      environmentPaths: userConfig.environmentPaths,
      logLineLimit: userConfig.logLineLimit,
      locale: i18n.language,
    },
  });

  const userFormWatcher = userForm.watch();

  useEffect(() => {
    const currentValues = userForm.getValues();

    if (lastUserSavedValues.current === null) {
      lastUserSavedValues.current = JSON.parse(JSON.stringify(currentValues));
      return;
    }

    if (
      JSON.stringify(currentValues) ===
      JSON.stringify(lastUserSavedValues.current)
    ) {
      return;
    }

    setHasUserPendingChanges(true);
    const timeout = setTimeout(async () => {
      lastUserSavedValues.current = JSON.parse(JSON.stringify(currentValues));
      await userForm.handleSubmit(saveUserSettingsForm)();
      setHasUserPendingChanges(false);
    }, 300);

    return () => clearTimeout(timeout);
  }, [userForm, userFormWatcher]);

  const [supportedLanguages, setSupportedLanguages] = useState<
    SupportedLanguage[]
  >([]);

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

  //endregion

  const hasFormErrors =
    Object.keys(projectForm.formState.errors).length > 0 ||
    Object.keys(userForm.formState.errors).length > 0;

  const value: SettingsContextData = {
    initialTab,
    hasPendingChanges:
      hasProjectPendingChanges || hasUserPendingChanges || hasFormErrors,
    projectSettingsForm: projectForm,
    userSettingsForm: userForm,
    supportedLanguages,
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
