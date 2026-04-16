import { z } from "zod";

export const formSchema = z.object({
  name: z.string().min(1, "projectForm.validation.nameRequired"),
  baseWorkingDirectory: z.string().min(1, "projectForm.validation.baseDirRequired"),
});
export type FormSchemaType = z.infer<typeof formSchema>;
