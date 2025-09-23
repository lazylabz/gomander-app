import { ChartNoAxesGantt, Save } from "lucide-react";

import { BaseWorkingDirectoryField } from "@/components/modals/Project/common/BaseWorkingDirectoryField.tsx";
import { ProjectNameField } from "@/components/modals/Project/common/ProjectNameField.tsx";
import { Button } from "@/design-system/components/ui/button.tsx";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/design-system/components/ui/card.tsx";
import { Form } from "@/design-system/components/ui/form.tsx";
import { useSettingsContext } from "@/screens/SettingsScreen/contexts/settingsContext.tsx";
import type { SettingsFormType } from "@/screens/SettingsScreen/contexts/settingsFormSchema.ts";
import { ProjectFailurePatternsField } from "./ProjectFailurePatterns";

export const ProjectSettings = () => {
  const { settingsForm, saveSettings, hasUnsavedChanges } =
    useSettingsContext();

  return (
    <Form {...settingsForm}>
      <form
        onSubmit={settingsForm.handleSubmit(saveSettings)}
        className="w-full h-full flex flex-col justify-between"
      >
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <ChartNoAxesGantt size={20} />
              <span>Project information</span>
            </CardTitle>
            <CardDescription>
              Configure your project details and basic settings.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col gap-4">
              <ProjectNameField<SettingsFormType> />
              <BaseWorkingDirectoryField<SettingsFormType> />
              <ProjectFailurePatternsField />  {/* New regex field */}
            </div>
          </CardContent>
        </Card>

        <Button
          className="self-end cursor-pointer"
          type="submit"
          disabled={!hasUnsavedChanges}
        >
          <Save />
          Save
        </Button>
      </form>
    </Form>
  );
};
