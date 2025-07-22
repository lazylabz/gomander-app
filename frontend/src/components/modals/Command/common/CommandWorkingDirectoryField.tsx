import { useFormContext } from "react-hook-form";

import type { FormSchemaType } from "@/components/modals/Command/common/formSchema.ts";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";

export const CommandWorkingDirectoryField = () => {
  const form = useFormContext<FormSchemaType>();

  return (
    <FormField
      control={form.control}
      name="workingDirectory"
      render={({ field }) => (
        <FormItem>
          <FormLabel>Command</FormLabel>
          <FormControl>
            <Input placeholder={"/Users/hackerman/Code"} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
