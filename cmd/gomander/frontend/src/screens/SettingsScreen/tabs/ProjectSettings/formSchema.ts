import { z } from "zod";

import { availableThemes } from "@/contexts/theme.tsx";

export const formSchema = z.object({
  environmentPaths: z.array(
    z.object({
      value: z.string().min(1, "Path cannot be empty"),
    }),
  ),
  theme: z.enum(availableThemes),
});

export type FormType = z.infer<typeof formSchema>;
