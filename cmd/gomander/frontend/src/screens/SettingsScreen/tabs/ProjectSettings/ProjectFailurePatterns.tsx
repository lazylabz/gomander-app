import type { FC } from "react";
import { useFormContext } from "react-hook-form";
import type { SettingsFormType } from "../../contexts/settingsFormSchema";

export const ProjectFailurePatternsField: FC = () => {
  const { register } = useFormContext<SettingsFormType>();

  return (
    <div className="mb-4">
      <label className="block font-medium mb-1">Failure Regex Patterns</label>
      <textarea
        {...register("failurePatterns")}
        placeholder="One regex per line, e.g., nodemon.*app crashed"
        className="w-full h-24 border rounded p-2"
      />
    </div>
  );
};
