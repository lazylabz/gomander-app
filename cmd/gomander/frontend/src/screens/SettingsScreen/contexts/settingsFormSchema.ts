import { z } from "zod";

import { availableThemes } from "@/contexts/theme.tsx";

export const settingsFormSchema = z.object({
  // User settings schema
  environmentPaths: z.array(
    z.object({
      value: z.string().min(1, "Path cannot be empty"),
    }),
  ),
  theme: z.enum(availableThemes),
  // Project settings schema
  name: z.string().min(1, "Project name is required"),
  baseWorkingDirectory: z.string().min(1, "Base working directory is required"),
});

export type SettingsFormSchemaType = z.infer<typeof settingsFormSchema>;
