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

export const CommandCommandField = () => {
  const form = useFormContext<FormSchemaType>();

  return (
    <FormField
      control={form.control}
      name="command"
      render={({ field }) => (
        <FormItem>
          <FormLabel>Command</FormLabel>
          <FormControl>
            <Input placeholder={'cowsay "Hello World!"'} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
