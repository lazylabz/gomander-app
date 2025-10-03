import { z } from "zod";

import { availableThemes } from "@/contexts/theme.tsx";

export const settingsFormSchema = z.object({
  // User settings schema
  environmentPaths: z.array(
    z.object({
      id: z.uuid(),
      path: z.string().min(1, "Path cannot be empty"),
    }),
  ),
  theme: z.enum(availableThemes),
  logLineLimit: z
    .number()
    .int()
    .min(1, "Must be at least 1")
    .max(5000, "Must be at most 5000"),
  // Project settings schema
  name: z.string().min(1, "Project name is required"),
  baseWorkingDirectory: z.string().min(1, "Base working directory is required"),
});

export type SettingsFormType = z.infer<typeof settingsFormSchema>;
