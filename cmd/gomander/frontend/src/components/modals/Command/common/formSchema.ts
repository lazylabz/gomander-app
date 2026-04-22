import { z } from "zod";

import i18n from "@/design-system/lib/i18n.ts";

export const formSchema = z.object({
	name: z.string().min(1, {
		error: () => i18n.t("commandForm.validation.nameRequired"),
	}),
	command: z.string().min(1, {
		error: () => i18n.t("commandForm.validation.commandRequired"),
	}),
	workingDirectory: z.string().min(0),
	link: z.string().min(0),
	errorPatterns: z.string().min(0),
});
export type FormSchemaType = z.infer<typeof formSchema>;
