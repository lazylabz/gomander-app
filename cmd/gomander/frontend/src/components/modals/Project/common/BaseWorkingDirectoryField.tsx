import { type Path, useFormContext } from "react-hook-form";

import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";

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
            <Input placeholder="/Users/hackerman/Code" {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
