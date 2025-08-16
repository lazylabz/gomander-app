import { z } from "zod";

export const formSchema = z.object({
  name: z.string().min(1, {
    message: "Command group name is required",
  }),
  commands: z.array(z.string()).min(1, {
    message: "You must add at least one command to the group",
  }),
});
export type FormSchemaType = z.infer<typeof formSchema>;
