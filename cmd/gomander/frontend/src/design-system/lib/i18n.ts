import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import { z } from 'zod';

import { dataService, translationsService } from '@/contracts/service';
import type { Localization } from '@/contracts/types';

// Type i18next to use our Localization interface
declare module 'i18next' {
  interface CustomTypeOptions {
    resources: {
      translation: Localization;
    };
  }
}

export const initI18n = async () => {
  // Get supported languages and user's current locale from backend
  const [supportedLanguages, userConfig] = await Promise.all([
    translationsService.getSupportedLanguages(),
    dataService.getUserConfig()
  ]);

  // Load translations for current language
  const currentLang = userConfig.locale || 'en';
  const translations = await translationsService.getTranslation(currentLang);

  await i18n
    .use(initReactI18next)
    .init({
      lng: currentLang,
      fallbackLng: 'en',
      supportedLngs: supportedLanguages,

      resources: {
        [currentLang]: {
          translation: translations
        }
      },

      ns: ['translation'],
      defaultNS: 'translation',

      interpolation: {
        escapeValue: false,
      },

      react: {
        useSuspense: false,
      }
    });

  // Translate Zod validation messages by treating the message string as an i18n key.
  // Schemas set their messages to i18n key strings (e.g. "commandForm.validation.nameRequired")
  // and this errorMap resolves them at validation time so language changes are reflected.
  z.setErrorMap((issue) => {
    const key = issue.message;
    if (key && i18n.exists(key)) {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      return { message: i18n.t(key as any) };
    }
    return { message: issue.message ?? 'Invalid value' };
  });

  return i18n;
};

export default i18n;
