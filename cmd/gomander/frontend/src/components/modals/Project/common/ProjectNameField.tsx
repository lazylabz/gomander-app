import { type Path, useFormContext } from "react-hook-form";

import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";

export const ProjectNameField = <T extends { name: string }>() => {
  const form = useFormContext<T>();

  const name = "name" satisfies keyof T as Path<T>;

  return (
    <FormField
      control={form.control}
      name={name}
      render={({ field }) => (
        <FormItem>
          <FormLabel>Name</FormLabel>
          <FormControl>
            <Input
              autoComplete="off"
              autoCorrect="off"
              autoCapitalize="off"
              placeholder="My awesome project"
              {...field}
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
