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

export const ProjectNameField = () => {
  const form = useFormContext<FormSchemaType>();

  return (
    <FormField
      control={form.control}
      name="name"
      render={({ field }) => (
        <FormItem>
          <FormLabel>Name</FormLabel>
          <FormControl>
            <Input placeholder="My awesome project" {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
