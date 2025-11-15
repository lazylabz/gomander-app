import { ChartNoAxesGantt } from "lucide-react";

import { BaseWorkingDirectoryField } from "@/components/modals/Project/common/BaseWorkingDirectoryField.tsx";
import { ProjectNameField } from "@/components/modals/Project/common/ProjectNameField.tsx";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/design-system/components/ui/card.tsx";
import { Form } from "@/design-system/components/ui/form.tsx";
import { useSettingsContext } from "@/screens/SettingsScreen/context/settingsContext.tsx";
import { type ProjectSettingsSchemaType } from "@/screens/SettingsScreen/schemas/projectSettingsSchema.ts";

export const ProjectSettings = () => {
  const { projectSettingsForm } = useSettingsContext();

  return (
    <Form {...projectSettingsForm}>
      <form className="w-full h-full flex flex-col justify-between">
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
              <ProjectNameField<ProjectSettingsSchemaType> />
              <BaseWorkingDirectoryField<ProjectSettingsSchemaType> />
            </div>
          </CardContent>
        </Card>
      </form>
    </Form>
  );
};
