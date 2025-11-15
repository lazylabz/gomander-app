import { zodResolver } from "@hookform/resolvers/zod";
import { ChartNoAxesGantt } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";

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
import {
  projectSettingsSchema,
  type ProjectSettingsSchemaType,
} from "@/screens/SettingsScreen/schemas/projectSettingsSchema.ts";
import { saveProjectSettingsForm } from "@/screens/SettingsScreen/useCases/saveProjectSettingsForm.ts";
import { useProjectStore } from "@/store/projectStore.ts";

export const ProjectSettings = () => {
  const projectInfo = useProjectStore((state) => state.projectInfo);

  const [isSaved, setIsSaved] = useState(true);
  const lastSavedValues = useRef<ProjectSettingsSchemaType | null>(null);

  const form = useForm<ProjectSettingsSchemaType>({
    resolver: zodResolver(projectSettingsSchema),
    defaultValues: {
      name: projectInfo?.name || "",
      baseWorkingDirectory: projectInfo?.workingDirectory || "",
    },
  });

  const formWatcher = form.watch();

  useEffect(() => {
    const currentValues = form.getValues();

    if (lastSavedValues.current === null) {
      lastSavedValues.current = JSON.parse(JSON.stringify(currentValues));
      return;
    }

    if (
      JSON.stringify(currentValues) === JSON.stringify(lastSavedValues.current)
    ) {
      return;
    }

    setIsSaved(false);
    const timeout = setTimeout(async () => {
      lastSavedValues.current = JSON.parse(JSON.stringify(currentValues));
      await form.handleSubmit(saveProjectSettingsForm)();
      setIsSaved(true);
    }, 300);

    return () => clearTimeout(timeout);
  }, [form, formWatcher]);

  return (
    <Form {...form}>
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
            <CardDescription>
              {/* TODO: Move this to page header */}
              <span>{isSaved ? "Saved" : "Saving..."}</span>
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
