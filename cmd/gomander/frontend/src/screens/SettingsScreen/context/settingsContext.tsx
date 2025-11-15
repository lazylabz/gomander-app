import { zodResolver } from "@hookform/resolvers/zod";
import { createContext, useContext, useEffect, useRef, useState } from "react";
import { useForm, type UseFormReturn } from "react-hook-form";
import { useLocation } from "react-router";

import {
  projectSettingsSchema,
  type ProjectSettingsSchemaType,
} from "@/screens/SettingsScreen/schemas/projectSettingsSchema.ts";
import { saveProjectSettingsForm } from "@/screens/SettingsScreen/useCases/saveProjectSettingsForm.ts";
import { useProjectStore } from "@/store/projectStore.ts";

export enum SettingsTab {
  User = "user",
  Project = "project",
}

// Define context
export interface SettingsContextData {
  initialTab: SettingsTab;
  hasPendingChanges: boolean;
  projectSettingsForm: UseFormReturn<ProjectSettingsSchemaType>;
}

export const settingsContext = createContext<SettingsContextData>({
  initialTab: SettingsTab.User,
  hasPendingChanges: false,
  projectSettingsForm: {} as UseFormReturn<ProjectSettingsSchemaType>,
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

  const value: SettingsContextData = {
    initialTab,
    hasPendingChanges: hasProjectPendingChanges,
    projectSettingsForm: projectForm,
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
