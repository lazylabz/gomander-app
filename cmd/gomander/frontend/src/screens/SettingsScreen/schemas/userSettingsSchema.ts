import { z } from "zod";

import { LOGS_FONT_SIZE } from "@/constants/logsFontSize.ts";
import i18n from "@/design-system/lib/i18n.ts";

const logsFontSizeOutOfRange = {
	error: () =>
		i18n.t("userSettingsForm.validation.logsFontSize", {
			min: LOGS_FONT_SIZE.min,
			max: LOGS_FONT_SIZE.max,
		}),
};

export const userSettingsSchema = z.object({
	environmentPaths: z.array(
		z.object({
			id: z.uuid(),
			path: z.string().min(1, {
				error: () => i18n.t("userSettingsForm.validation.pathEmpty"),
			}),
		}),
	),
	locale: z.string(),
	logsFontSize: z
		.number(logsFontSizeOutOfRange)
		.int(logsFontSizeOutOfRange)
		.min(LOGS_FONT_SIZE.min, logsFontSizeOutOfRange)
		.max(LOGS_FONT_SIZE.max, logsFontSizeOutOfRange),
});

export type UserSettingsSchemaType = z.infer<typeof userSettingsSchema>;
