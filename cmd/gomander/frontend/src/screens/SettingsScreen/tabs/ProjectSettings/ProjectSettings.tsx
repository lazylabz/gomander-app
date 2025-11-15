import { zodResolver } from "@hookform/resolvers/zod";
import { ChartNoAxesGantt } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

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
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchProject } from "@/queries/fetchProject.ts";
import { projectStore, useProjectStore } from "@/store/projectStore.ts";
import { editOpenedProject } from "@/useCases/project/editOpenedProject.ts";

const formSchema = z.object({
  name: z.string().min(1, "Project name is required"),
  baseWorkingDirectory: z.string().min(1, "Base working directory is required"),
});

type FormSchemaType = z.infer<typeof formSchema>;

const handleSave = async (formData: FormSchemaType) => {
  const { projectInfo } = projectStore.getState();
  if (!projectInfo) {
    return;
  }

  try {
    await editOpenedProject({
      ...projectInfo,
      name: formData.name,
      workingDirectory: formData.baseWorkingDirectory,
    });
    toast.success("Project settings saved successfully");
  } catch (e) {
    toast.error(parseError(e, "Failed to save project settings"));
  }

  await fetchProject();
};

export const ProjectSettings = () => {
  const projectInfo = useProjectStore((state) => state.projectInfo);

  const [isSaved, setIsSaved] = useState(true);
  const lastSavedValues = useRef<FormSchemaType | null>(null);

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
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
      await form.handleSubmit(handleSave)();
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
              <ProjectNameField<FormSchemaType> />
              <BaseWorkingDirectoryField<FormSchemaType> />
            </div>
          </CardContent>
        </Card>
      </form>
    </Form>
  );
};
