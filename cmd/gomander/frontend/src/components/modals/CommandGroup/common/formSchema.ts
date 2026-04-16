import { z } from "zod";

export const formSchema = z.object({
  name: z.string().min(1, {
    message: "commandGroupForm.validation.nameRequired",
  }),
  commands: z.array(z.string()).min(1, {
    message: "commandGroupForm.validation.commandsRequired",
  }),
});
export type FormSchemaType = z.infer<typeof formSchema>;
