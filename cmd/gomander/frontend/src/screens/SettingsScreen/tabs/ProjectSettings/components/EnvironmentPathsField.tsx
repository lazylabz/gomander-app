import { Plus, Trash } from "lucide-react";
import { useFieldArray, useFormContext } from "react-hook-form";

import { Button } from "@/components/ui/button.tsx";
import {
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import type { FormType } from "@/screens/SettingsScreen/tabs/ProjectSettings/formSchema.ts";

export const EnvironmentPathsField = () => {
  const { control } = useFormContext<FormType>();

  const { fields, append, remove } = useFieldArray({
    control,
    name: "environmentPaths" as const,
  });

  const addNewPath = () => {
    append({ value: "" });
  };

  const removePath = (index: number) => () => {
    remove(index);
  };
  return (
    <div className="flex flex-col items-center">
      {fields.length !== 0 && (
        <div className="flex flex-col gap-2 w-full">
          {fields.map((_, index) => (
            <div
              className="flex flex-row w-full items-center gap-3 justify-between"
              key={index}
            >
              <FormField
                control={control}
                name={`environmentPaths.${index}.value` as const}
                render={({ field }) => (
                  <FormItem className="flex-1">
                    <FormControl>
                      <Input
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
