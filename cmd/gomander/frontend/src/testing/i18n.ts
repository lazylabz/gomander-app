import { initReactI18next } from "react-i18next";
import { onTestFinished } from "vitest";

import i18n from "@/design-system/lib/i18n.ts";

// No resources are registered, so i18next echoes every key back: a test asserts
// on the key instead of on the English copy. Registering with react-i18next
// makes rendered components translate through this same instance.
export const installTranslations = async (): Promise<void> => {
	if (i18n.isInitialized) {
		return;
	}

	await i18n.use(initReactI18next).init({ lng: "en", resources: {} });
};

// The echo drops interpolated values, so a test that needs to see one gives its
// key real copy for the length of that test. i18next cannot drop a single key,
// so the whole bundle goes; it only ever holds copy a test put there.
export const withTranslation = (key: string, copy: string): void => {
	i18n.addResource("en", "translation", key, copy);
	onTestFinished(() => {
		i18n.removeResourceBundle("en", "translation");
	});
};
