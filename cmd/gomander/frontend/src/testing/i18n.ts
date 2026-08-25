import i18n from "@/design-system/lib/i18n.ts";

// No resources are registered, so i18next echoes every key back: a test asserts
// on the key instead of on the English copy.
export const installTranslations = async (): Promise<void> => {
	if (i18n.isInitialized) {
		return;
	}

	await i18n.init({ lng: "en", resources: {} });
};
