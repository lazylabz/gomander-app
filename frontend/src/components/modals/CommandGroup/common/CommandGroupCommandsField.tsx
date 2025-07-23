import { useFormContext } from "react-hook-form";

import type { FormSchemaType } from "@/components/modals/CommandGroup/common/formSchema.ts";
import { Checkbox } from "@/components/ui/checkbox.tsx";
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import { useDataContext } from "@/contexts/DataContext.tsx";

export const CommandGroupCommandsField = () => {
  const { commands } = useDataContext();

  const form = useFormContext<FormSchemaType>();

  return (
    <FormField
      control={form.control}
      name="commands"
      render={() => (
        <FormItem>
          <div className="mb-2">
            <FormLabel className="text-base">Commands</FormLabel>
            <FormDescription>
              Don't worry about the order, you'll be able to change it later
            </FormDescription>
          </div>
          {Object.values(commands).map((item) => (
            <FormField
              key={item.id}
              control={form.control}
              name="commands"
              render={({ field }) => {
                return (
                  <FormItem
                    key={item.id}
                    className="flex flex-row items-center gap-2"
                  >
                    <FormControl>
                      <Checkbox
                        checked={field.value?.includes(item.id)}
                        onCheckedChange={(checked) => {
                          return checked
                            ? field.onChange([...field.value, item.id])
                            : field.onChange(
                                field.value?.filter(
                                  (value) => value !== item.id,
                                ),
                              );
                        }}
                      />
                    </FormControl>
                    <FormLabel className="text-sm font-normal">
                      {item.name}
                    </FormLabel>
                  </FormItem>
                );
              }}
            />
          ))}
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
