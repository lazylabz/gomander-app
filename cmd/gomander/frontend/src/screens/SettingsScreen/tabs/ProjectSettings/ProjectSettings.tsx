import { zodResolver } from "@hookform/resolvers/zod";
import { ChartNoAxesGantt, Save } from "lucide-react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";

import { BaseWorkingDirectoryField } from "@/components/modals/Project/common/BaseWorkingDirectoryField.tsx";
import {
  formSchema,
  type FormSchemaType,
} from "@/components/modals/Project/common/createAndEditSchema.ts";
import { ProjectNameField } from "@/components/modals/Project/common/ProjectNameField.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card.tsx";
import { Form } from "@/components/ui/form.tsx";
import { fetchProject } from "@/queries/fetchProject.ts";
import { useProjectStore } from "@/store/projectStore.ts";
import { editOpenedProject } from "@/useCases/project/editOpenedProject.ts";

export const ProjectSettings = () => {
  const projectInfo = useProjectStore((state) => state.projectInfo);

  const navigate = useNavigate();

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    values: {
      name: projectInfo?.name || "",
      baseWorkingDirectory: projectInfo?.baseWorkingDirectory || "",
    },
  });
  const onSubmit = async (values: FormSchemaType) => {
    if (!projectInfo) {
      return;
    }

    await editOpenedProject({
      ...projectInfo,
      name: values.name,
      baseWorkingDirectory: values.baseWorkingDirectory,
    });

    await fetchProject();

    navigate(-1);
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
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
              <ProjectNameField<FormSchemaType> />
              <BaseWorkingDirectoryField<FormSchemaType> />
            </div>
          </CardContent>
        </Card>

        <Button className="self-end" type="submit">
          <Save />
          Save
        </Button>
      </form>
    </Form>
  );
};
