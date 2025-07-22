import { z } from "zod";

export const formSchema = z.object({
  name: z.string().min(1, {
    message: "Command name is required",
  }),
  command: z.string().min(1, {
    message: "Command is required",
  }),
  workingDirectory: z.string().min(0),
});
export type FormSchemaType = z.infer<typeof formSchema>;
