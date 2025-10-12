import type { FC } from "react";
import { useState } from "react";
import { Controller,useFormContext } from "react-hook-form";

import type { SettingsFormType } from "../../contexts/settingsFormSchema";



export const ProjectFailurePatternsField: FC = () => {
  const { control } = useFormContext<SettingsFormType>();

  const [displayValue, setDisplayValue] = useState("");

  return (
    <div className="mb-4">
      <label className="block font-medium mb-1">Failure Regex Patterns</label>
      <Controller
        name="failurePatterns"
        control={control}
        render={({ field }) => {
          if (!displayValue && field.value) {
            setDisplayValue((field.value || []).join("\n"));
          }

          return (
            <textarea
              value={displayValue}
              onChange={(e) => {
                const newValue = e.target.value;
                setDisplayValue(newValue);
                
                const lines = newValue.split("\n").map((line) => line.trim()).filter(Boolean);
                field.onChange(lines);
              }}
              placeholder="One regex per line, e.g., nodemon.*app crashed"
              className="w-full h-24 border rounded p-2"
            />
          );
        }}
      />
    </div>
  );
};