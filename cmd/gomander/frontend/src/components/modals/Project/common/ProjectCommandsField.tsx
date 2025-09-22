import {useFormContext} from "react-hook-form";

import type {FormSchemaType} from "@/components/modals/Project/common/importAndExportSchema.ts";
import {Checkbox} from "@/components/ui/checkbox.tsx";
import {FormControl, FormField, FormItem, FormLabel,} from "@/components/ui/form.tsx";
import type {ProjectExport} from "@/contracts/types.ts";

export const ProjectCommandsField = ({
  commands,
}: {
  commands: ProjectExport["commands"];
}) => {
  const form = useFormContext<FormSchemaType>();

  return (
    <FormItem className="flex-1">
      <FormLabel className="mb-1">Commands</FormLabel>
      <div className="max-h-[300px] flex flex-col gap-2 overflow-y-auto pr-2">
        {commands.map((command) => (
          <FormField
            key={command.id}
            control={form.control}
            name="commands"
            render={({ field }) => (
              <FormItem key={command.id} className="flex flex-row items-center gap-2">
                <FormControl>
                  <Checkbox
                    checked={field.value?.includes(command.id)}
                    onCheckedChange={(checked) => {
                      return checked
                        ? field.onChange([...field.value, command.id])
                        : field.onChange(
                            field.value?.filter((value) => value !== command.id),
                          );
                    }}
                  />
                </FormControl>
                <FormLabel
                  title={command.name}
                  className="text-sm font-normal truncate"
                >
                  {command.name}
                </FormLabel>
              </FormItem>
            )}
          />
        ))}
      </div>
    </FormItem>
  );
};
