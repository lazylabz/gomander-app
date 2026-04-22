import { z } from "zod";

import i18n from "@/design-system/lib/i18n.ts";

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
	logLineLimit: z
		.number()
		.int()
		.min(1, { error: () => i18n.t("userSettingsForm.validation.logLimitMin") })
		.max(5000, {
			error: () => i18n.t("userSettingsForm.validation.logLimitMax"),
		}),
});

export type UserSettingsSchemaType = z.infer<typeof userSettingsSchema>;
