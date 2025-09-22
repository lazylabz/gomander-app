import {useFormContext} from "react-hook-form";

import type {FormSchemaType} from "@/components/modals/Project/common/importAndExportSchema.ts";
import {Checkbox} from "@/components/ui/checkbox.tsx";
import {FormControl, FormField, FormItem, FormLabel,} from "@/components/ui/form.tsx";
import {Tooltip, TooltipContent, TooltipTrigger} from "@/components/ui/tooltip.tsx";
import type {ProjectExport} from "@/contracts/types.ts";

export const ProjectCommandGroupsField = ({
  commandGroups,
  selectedCommandIds,
                                          }: {
  commandGroups: ProjectExport["commandGroups"],
  selectedCommandIds: string[]
}) => {
  const form = useFormContext<FormSchemaType>()

  return <FormItem className="flex-1">
    <FormLabel className="mb-1">Command Groups</FormLabel>
    <div className="max-h-[300px] flex flex-col gap-2 overflow-y-auto pr-2">
      {commandGroups.map((commandGroup) => {
        return (
          <FormField
            key={commandGroup.id}
            control={form.control}
            name="commandGroups"
            render={({ field }) => {
              const disabled = !commandGroup.commandIds.some((id) =>
                selectedCommandIds.includes(id),
              );

              return (
                <FormItem
                  key={commandGroup.id}
                  className="flex flex-row items-center gap-2"
                >
                  <FormControl>
                    <Checkbox
                      disabled={disabled}
                      checked={field.value?.includes(commandGroup.id)}
                      onCheckedChange={(checked) => {
                        return checked
                          ? field.onChange([
                            ...field.value,
                            commandGroup.id,
                          ])
                          : field.onChange(
                            field.value?.filter(
                              (value) => value !== commandGroup.id,
                            ),
                          );
                      }}
                    />
                  </FormControl>

                  <Tooltip>
                    <TooltipTrigger asChild>
                      <FormLabel
                        title={commandGroup.name}
                        className="text-sm font-normal truncate"
                      >
                        {commandGroup.name}
                      </FormLabel>
                    </TooltipTrigger>
                    {disabled && (
                      <TooltipContent>
                        Select at least one of its commands to
                        enable this group
                      </TooltipContent>
                    )}
                  </Tooltip>
                </FormItem>
              );
            }}
          />
        );
      })}
    </div>
  </FormItem>
}