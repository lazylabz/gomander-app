import { z } from "zod";

import i18n from "@/design-system/lib/i18n.ts";

export const projectSettingsSchema = z.object({
	name: z
		.string()
		.min(1, { error: () => i18n.t("projectForm.validation.nameRequired") }),
	baseWorkingDirectory: z
		.string()
		.min(1, { error: () => i18n.t("projectForm.validation.baseDirRequired") }),
});

export type ProjectSettingsSchemaType = z.infer<typeof projectSettingsSchema>;
