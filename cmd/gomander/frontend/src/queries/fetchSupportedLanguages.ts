import { LANGUAGE_LABELS } from "@/constants/languages.ts";
import { translationsService } from "@/contracts/service.ts";

export type SupportedLanguage = {
	value: string;
	label: string;
};

export const fetchSupportedLanguages = async (): Promise<
	SupportedLanguage[]
> => {
	const languages = await translationsService.getSupportedLanguages();

	return languages.map((lang) => ({
		value: lang,
		label: LANGUAGE_LABELS[lang] || lang,
	}));
};
