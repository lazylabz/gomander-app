import { Plus, Trash } from "lucide-react";
import { useFieldArray, useFormContext } from "react-hook-form";

import { FSInput } from "@/components/inputs/FSInput.tsx";
import { Button } from "@/design-system/components/ui/button.tsx";
import {
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/design-system/components/ui/form.tsx";
import type { SettingsFormType } from "@/screens/SettingsScreen/contexts/settingsFormSchema.ts";

export const EnvironmentPathsField = () => {
  const { control } = useFormContext<SettingsFormType>();

  const { fields, append, remove } = useFieldArray({
    control,
    name: "environmentPaths" as const,
  });

  const addNewPath = () => {
    append({ id: crypto.randomUUID(), path: "" });
  };

  const removePath = (index: number) => () => {
    remove(index);
  };

  return (
    <div className="flex flex-col items-center">
      {fields.length !== 0 && (
        <div className="flex flex-col gap-2 w-full">
          {fields.map((field, index) => (
            <div
              className="flex flex-row w-full items-center gap-3 justify-between"
              key={field.id}
            >
              <FormField
                control={control}
                name={`environmentPaths.${index}.path` as const}
                render={({ field }) => (
                  <FormItem className="flex-1">
                    <FormControl>
                      <FSInput
                        autoComplete="off"
                        autoCorrect="off"
                        autoCapitalize="off"
                        placeholder="/Users/hackerman/.nvm/versions/node/v20.18.0/bin"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Trash
                size={18}
                className="text-muted-foreground cursor-pointer hover:text-destructive"
                onClick={removePath(index)}
              />
            </div>
          ))}
        </div>
      )}
      <Button
        size="sm"
        type="button"
        variant="ghost"
        className="mt-2 w-auto"
        onClick={addNewPath}
      >
        Add
        <Plus />
      </Button>
    </div>
  );
};
