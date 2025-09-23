import type { FC } from "react";
import { useFormContext, Controller } from "react-hook-form";
import type { SettingsFormType } from "../../contexts/settingsFormSchema";

export const ProjectFailurePatternsField: FC = () => {
  const { control } = useFormContext<SettingsFormType>();

  return (
    <div className="mb-4">
      <label className="block font-medium mb-1">Failure Regex Patterns</label>
      <Controller
        name="failurePatterns"
        control={control}
        render={({ field }) => (
          <textarea
            value={(field.value || []).join("\n")}
            onChange={(e) =>
              field.onChange(e.target.value.split("\n").map((line) => line.trim()).filter(Boolean))
            }
            placeholder="One regex per line, e.g., nodemon.*app crashed"
            className="w-full h-24 border rounded p-2"
          />
        )}
      />
    </div>
  );
};