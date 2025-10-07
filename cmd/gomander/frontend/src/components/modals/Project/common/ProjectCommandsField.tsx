import { useFormContext } from "react-hook-form";

import type { FormSchemaType } from "@/components/modals/Project/common/importAndExportSchema.ts";
import type { ProjectExport } from "@/contracts/types.ts";
import { Checkbox } from "@/design-system/components/ui/checkbox.tsx";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from "@/design-system/components/ui/form.tsx";

export const ProjectCommandsField = ({
  commands,
  onChange,
}: {
  commands: ProjectExport["commands"];
  onChange?: (selectedCommandIds: string[]) => void;
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
              <FormItem
                key={command.id}
                className="flex flex-row items-center gap-2"
              >
                <FormControl>
                  <Checkbox
                    checked={field.value?.includes(command.id)}
                    onCheckedChange={(checked) => {
                      const newValue = checked
                        ? [...field.value, command.id]
                        : field.value?.filter((value) => value !== command.id);

                      onChange?.(newValue);

                      return field.onChange(newValue);
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
