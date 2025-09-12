import { useFormContext } from "react-hook-form";

import { FSInput } from "@/components/inputs/FSInput.tsx";
import type { FormSchemaType } from "@/components/modals/Command/common/formSchema.ts";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";

export const CommandWorkingDirectoryField = () => {
  const form = useFormContext<FormSchemaType>();

  return (
    <FormField
      control={form.control}
      name="workingDirectory"
      render={({ field }) => (
        <FormItem>
          <FormLabel>Working Directory</FormLabel>
          <FormControl>
            <FSInput
              autoComplete="off"
              autoCorrect="off"
              autoCapitalize="off"
              {...field}
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
