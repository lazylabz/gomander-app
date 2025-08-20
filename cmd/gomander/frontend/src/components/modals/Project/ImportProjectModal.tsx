import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { toast } from "sonner";

import { BaseWorkingDirectoryField } from "@/components/modals/Project/common/BaseWorkingDirectoryField.tsx";
import {
  formSchema,
  type FormSchemaType,
} from "@/components/modals/Project/common/importAndExportSchema.ts";
import { ProjectNameField } from "@/components/modals/Project/common/ProjectNameField.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog.tsx";
import { Form } from "@/components/ui/form.tsx";
import type { ProjectExport } from "@/contracts/types.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { importProject } from "@/useCases/project/importProject.ts";

export const ImportProjectModal = ({
  open,
  onSuccess,
  onClose,
  project,
}: {
  open: boolean;
  onSuccess: () => Promise<void>;
  onClose: () => void;
  project: ProjectExport | null;
}) => {
  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    values: {
      name: project?.name || "",
      baseWorkingDirectory: "",
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

    try {
      await importProject(project, values.name, values.baseWorkingDirectory);

      await onSuccess();
      handleOpenChange(false);
      toast.success("Project imported successfully");
    } catch (e) {
      toast.error("Failed to import the project: " + parseError(e));
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Import project</DialogTitle>
          <DialogDescription>
            Feel free to modify the values to the ones you prefer
          </DialogDescription>
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
