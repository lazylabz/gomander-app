import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { BaseWorkingDirectoryField } from "@/components/modals/Project/common/BaseWorkingDirectoryField.tsx";
import {
  formSchema,
  type FormSchemaType,
} from "@/components/modals/Project/common/createAndEditSchema.ts";
import { ProjectNameField } from "@/components/modals/Project/common/ProjectNameField.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Dialog, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { DialogContent } from "@/components/ui/dialog.tsx";
import { Form } from "@/components/ui/form.tsx";
import type { ProjectInfo } from "@/contracts/types.ts";
import { editOpenedProject } from "@/useCases/project/editOpenedProject.ts";

export const EditOpenedProjectModal = ({
  project,
  open,
  onClose,
  onSuccess,
}: {
  project: ProjectInfo | null;
  open: boolean;
  onClose: () => void;
  onSuccess: () => Promise<void>;
}) => {
  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    values: {
      name: project?.name || "",
      baseWorkingDirectory: project?.baseWorkingDirectory || "",
    },
  });

  const handleOpenChange = (open: boolean) => {
    if (!open) {
      onClose();
      form.reset();
    }
  };

  const onSubmit = async (values: FormSchemaType) => {
    if (!project) {
      return;
    }

    await editOpenedProject({
      ...project,
      name: values.name,
      baseWorkingDirectory: values.baseWorkingDirectory,
    });

    onSuccess();
    handleOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit project</DialogTitle>
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmit)}
              className="w-full mt-2 flex flex-col gap-4"
            >
              <ProjectNameField<FormSchemaType> />
              <BaseWorkingDirectoryField<FormSchemaType> />
              <Button className="self-end" type="submit">
                Save
              </Button>
            </form>
          </Form>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};
