import { type Path, useFormContext } from "react-hook-form";

import { FSInput } from "@/components/inputs/FSInput.tsx";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/design-system/components/ui/form.tsx";

export const BaseWorkingDirectoryField = <
  T extends { baseWorkingDirectory: string },
>() => {
  const form = useFormContext<T>();

  const name = "baseWorkingDirectory" satisfies keyof T as Path<T>;

  return (
    <FormField
      control={form.control}
      name={name}
      render={({ field }) => (
        <FormItem>
          <FormLabel>Base Working Directory</FormLabel>
          <FormControl>
            <FSInput
              autoComplete="off"
              autoCorrect="off"
              autoCapitalize="off"
              placeholder="/Users/hackerman/Code"
              {...field}
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
