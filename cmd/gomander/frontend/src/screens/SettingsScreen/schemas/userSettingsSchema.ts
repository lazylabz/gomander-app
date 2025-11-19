import { z } from "zod";

export const userSettingsSchema = z.object({
  environmentPaths: z.array(
    z.object({
      id: z.uuid(),
      path: z.string().min(1, "Path cannot be empty"),
    }),
  ),
  locale: z.string(),
  logLineLimit: z
    .number()
    .int()
    .min(1, "Must be at least 1")
    .max(5000, "Must be at most 5000"),
});

export type UserSettingsSchemaType = z.infer<typeof userSettingsSchema>;
