import { z } from "zod";

export const projectSettingsSchema = z.object({
  name: z.string().min(1, "projectForm.validation.nameRequired"),
  baseWorkingDirectory: z.string().min(1, "projectForm.validation.baseDirRequired"),
});

export type ProjectSettingsSchemaType = z.infer<typeof projectSettingsSchema>;
