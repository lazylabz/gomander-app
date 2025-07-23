import { z } from "zod";

export const formSchema = z.object({
  name: z.string().min(1, {
    message: "Command group name is required",
  }),
  commands: z.array(z.string()).min(1, {
    message: "You have to select at least one item.",
  }),
});
export type FormSchemaType = z.infer<typeof formSchema>;
