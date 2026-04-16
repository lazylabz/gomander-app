import { z } from "zod";

export const userSettingsSchema = z.object({
  environmentPaths: z.array(
    z.object({
      id: z.uuid(),
      path: z.string().min(1, "userSettingsForm.validation.pathEmpty"),
    }),
  ),
  locale: z.string(),
  logLineLimit: z
    .number()
    .int()
    .min(1, "userSettingsForm.validation.logLimitMin")
    .max(5000, "userSettingsForm.validation.logLimitMax"),
});

export type UserSettingsSchemaType = z.infer<typeof userSettingsSchema>;
