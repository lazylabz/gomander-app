import { ChartNoAxesGantt } from "lucide-react";
import { useTranslation } from "react-i18next";

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
  const { t } = useTranslation();
  const { projectSettingsForm } = useSettingsContext();

  return (
    <Form {...projectSettingsForm}>
      <form className="w-full h-full flex flex-col justify-between">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <ChartNoAxesGantt size={20} />
              <span>{t('projectSettingsForm.sectionTitle')}</span>
            </CardTitle>
            <CardDescription>
              {t('projectSettingsForm.sectionDescription')}
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
