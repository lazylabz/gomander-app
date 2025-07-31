import { useFormContext } from "react-hook-form";

import type { FormSchemaType } from "@/components/modals/Project/common/schema.ts";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";

export const BaseWorkingDirectoryField = () => {
  const form = useFormContext<FormSchemaType>();

  return (
    <FormField
      control={form.control}
      name="baseWorkingDirectory"
      render={({ field }) => (
        <FormItem>
          <FormLabel>Base Working Directory</FormLabel>
          <FormControl>
            <Input placeholder="/Users/hackerman/Code" {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
