import { z } from "zod";

export const formSchema = z.object({
  name: z.string().min(1, {
    message: "commandForm.validation.nameRequired",
  }),
  command: z.string().min(1, {
    message: "commandForm.validation.commandRequired",
  }),
  workingDirectory: z.string().min(0),
  link: z.string().min(0),
  errorPatterns: z.string().min(0),
});
export type FormSchemaType = z.infer<typeof formSchema>;
