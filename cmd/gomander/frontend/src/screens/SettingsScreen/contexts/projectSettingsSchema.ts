import { z } from "zod";

export const projectSettingsSchema = z.object({
  name: z.string().min(1, "Project name is required"),
  baseWorkingDirectory: z.string().min(1, "Base working directory is required"),
});

export type ProjectSettingsSchemaType = z.infer<typeof projectSettingsSchema>;
