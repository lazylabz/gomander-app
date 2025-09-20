import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import { dataService, translationsService } from '@/contracts/service';

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

  return i18n;
};

export default i18n;