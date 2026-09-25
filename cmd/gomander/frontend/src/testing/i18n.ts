import { initReactI18next } from "react-i18next";

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
