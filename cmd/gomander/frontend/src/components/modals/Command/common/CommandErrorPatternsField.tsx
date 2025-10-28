import { useFormContext } from "react-hook-form";

import type { FormSchemaType } from "@/components/modals/Command/common/formSchema.ts";
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/design-system/components/ui/form.tsx";
import { Textarea } from "@/design-system/components/ui/textarea.tsx";

export const CommandErrorPatternsField = () => {
  const form = useFormContext<FormSchemaType>();

  return (
    <FormField
      control={form.control}
      name="errorPatterns"
      render={({ field }) => (
        <FormItem>
          <FormLabel>Error patterns</FormLabel>
          <FormDescription className="text-xs">
            Regular expressions to identify error messages in the command
            output. Separate multiple patterns with new lines.
          </FormDescription>
          <FormControl>
            <Textarea
              autoComplete="off"
              autoCorrect="off"
              autoCapitalize="off"
              placeholder="/[nodemon] app crashed/gi"
              {...field}
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
