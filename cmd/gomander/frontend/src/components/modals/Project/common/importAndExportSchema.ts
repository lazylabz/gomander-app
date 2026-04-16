import { z } from "zod";

export const formSchema = z.object({
  name: z.string().min(1, "projectForm.validation.nameRequired"),
  baseWorkingDirectory: z.string().min(1, "projectForm.validation.baseDirRequired"),
  commands: z.array(z.string()),
  commandGroups: z.array(z.string()),
});
export type FormSchemaType = z.infer<typeof formSchema>;
