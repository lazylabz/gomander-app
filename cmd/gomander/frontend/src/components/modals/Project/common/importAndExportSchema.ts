import { z } from "zod";

export const formSchema = z.object({
  name: z.string().min(1, "Project name is required"),
  baseWorkingDirectory: z.string().min(1, "Base working directory is required"),
  commands: z.array(z.string()),
  commandGroups: z.array(z.string()),
});
export type FormSchemaType = z.infer<typeof formSchema>;
