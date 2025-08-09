import { z } from "zod";

export const formSchema = z.object({
  environmentPaths: z.array(
    z.object({
      value: z.string().min(1, "Path cannot be empty"),
    }),
  ),
});
export type FormType = z.infer<typeof formSchema>;
