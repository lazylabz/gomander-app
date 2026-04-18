import { z } from "zod";

import i18n from "@/design-system/lib/i18n.ts";

export const formSchema = z.object({
  name: z.string().min(1, { error: () => i18n.t("projectForm.validation.nameRequired") }),
  baseWorkingDirectory: z.string().min(1, { error: () => i18n.t("projectForm.validation.baseDirRequired") }),
  commands: z.array(z.string()),
  commandGroups: z.array(z.string()),
});
export type FormSchemaType = z.infer<typeof formSchema>;
