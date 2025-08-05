import { z } from "zod";

export const formSchema = z.object({
  name: z.string().min(1, "Project name is required"),
  baseWorkingDirectory: z.string().min(1, "Base working directory is required"),
});
export type FormSchemaType = z.infer<typeof formSchema>;
