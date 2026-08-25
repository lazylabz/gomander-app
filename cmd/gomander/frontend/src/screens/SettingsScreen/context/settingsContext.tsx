import { createContext, useContext, useEffect, useState } from "react";
import type { UseFormReturn } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { useLocation } from "react-router";

import { LANGUAGE_LABELS } from "@/constants/languages.ts";
import { translationsService } from "@/contracts/service.ts";
import { useAutosavedForm } from "@/hooks/useAutosavedForm.ts";
import {
	type ProjectSettingsSchemaType,
	projectSettingsSchema,
} from "@/screens/SettingsScreen/schemas/projectSettingsSchema.ts";
import {
	type UserSettingsSchemaType,
	userSettingsSchema,
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

	const projectInfo = useProjectStore((state) => state.projectInfo);
	const userConfig = useUserConfigurationStore((state) => state.userConfig);
	const { i18n } = useTranslation();

	const { form: projectForm, isPending: isProjectSavePending } =
		useAutosavedForm({
			schema: projectSettingsSchema,
			defaultValues: {
				name: projectInfo?.name || "",
				baseWorkingDirectory: projectInfo?.workingDirectory || "",
			},
			save: saveProjectSettingsForm,
		});

	const { form: userForm, isPending: isUserSavePending } = useAutosavedForm({
		schema: userSettingsSchema,
		defaultValues: {
			environmentPaths: userConfig.environmentPaths,
			locale: i18n.language,
		},
		save: saveUserSettingsForm,
	});

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
						label: LANGUAGE_LABELS[lang] || lang,
					}),
				);
				setSupportedLanguages(languageOptions);
			} catch (error) {
				console.error("Failed to load supported languages:", error);
			}
		};

		loadSupportedLanguages();
	}, []);

	const hasFormErrors =
		Object.keys(projectForm.formState.errors).length > 0 ||
		Object.keys(userForm.formState.errors).length > 0;

	const value: SettingsContextData = {
		initialTab,
		hasPendingChanges:
			isProjectSavePending || isUserSavePending || hasFormErrors,
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
