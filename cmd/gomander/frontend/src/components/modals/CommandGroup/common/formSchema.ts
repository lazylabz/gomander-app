import { z } from "zod";

import i18n from "@/design-system/lib/i18n.ts";

export const formSchema = z.object({
  name: z.string().min(1, {
    error: () => i18n.t("commandGroupForm.validation.nameRequired"),
  }),
  commands: z.array(z.string()).min(1, {
    error: () => i18n.t("commandGroupForm.validation.commandsRequired"),
  }),
});
export type FormSchemaType = z.infer<typeof formSchema>;
