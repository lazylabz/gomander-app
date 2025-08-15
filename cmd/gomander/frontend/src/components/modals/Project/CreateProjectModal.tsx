import { zodResolver } from "@hookform/resolvers/zod";
import type { SetStateAction } from "react";
import { useForm } from "react-hook-form";

import { BaseWorkingDirectoryField } from "@/components/modals/Project/common/BaseWorkingDirectoryField.tsx";
import {
  formSchema,
  type FormSchemaType,
} from "@/components/modals/Project/common/createSchema.ts";
import { ProjectNameField } from "@/components/modals/Project/common/ProjectNameField.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Dialog, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { DialogContent } from "@/components/ui/dialog.tsx";
import { Form } from "@/components/ui/form.tsx";
import { dataService } from "@/contracts/service";

export const CreateProjectModal = ({
  open,
  setOpen,
  onSuccess,
}: {
  open: boolean;
  setOpen: React.Dispatch<SetStateAction<boolean>>;
  onSuccess: () => Promise<void>;
}) => {
  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      baseWorkingDirectory: "",
    },
  });

  const handleOpenChange = (open: boolean) => {
    setOpen(open);
    if (!open) {
      form.reset();
    }
  };

  const onSubmit = async (values: FormSchemaType) => {
    await dataService.createProject({
      id: crypto.randomUUID(),
      name: values.name,
      workingDirectory: values.baseWorkingDirectory,
    });

    onSuccess();
    handleOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create project</DialogTitle>
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
