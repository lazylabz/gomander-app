import { getI18n } from "react-i18next";
import { toast } from "sonner";

import { translationsService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchUserConfig } from "@/queries/fetchUserConfig.ts";
import type { UserSettingsSchemaType } from "@/screens/SettingsScreen/schemas/userSettingsSchema.ts";
import { userConfigurationStore } from "@/store/userConfigurationStore.ts";
import { saveUserConfig } from "@/useCases/userConfig/saveUserConfig.ts";

const changeLanguage = async (lang: string) => {
  const i18nInstance = getI18n();

  if (i18nInstance.language === lang) {
    return;
  }

  if (!i18nInstance.hasResourceBundle(lang, "translation")) {
    const translations = await translationsService.getTranslation(lang);
    i18nInstance.addResourceBundle(lang, "translation", translations);
  }

  await i18nInstance.changeLanguage(lang);
};

export const saveUserSettingsForm = async (formData: UserSettingsSchemaType) => {
  const { userConfig } = userConfigurationStore.getState();

  try {
    await changeLanguage(formData.locale);
    await saveUserConfig({
      lastOpenedProjectId: userConfig.lastOpenedProjectId,
      environmentPaths: formData.environmentPaths,
      logLineLimit: formData.logLineLimit,
      locale: formData.locale,
    });
    toast.success(i18n.t('toast.settings.userSaveSuccess'));
  } catch (e) {
    toast.error(parseError(e, i18n.t('toast.settings.userSaveFailed')));
  }

  await fetchUserConfig();
};
