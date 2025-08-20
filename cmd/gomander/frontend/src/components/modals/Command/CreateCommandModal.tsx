import { zodResolver } from "@hookform/resolvers/zod";
import { Terminal } from "lucide-react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";

import { CommandCommandField } from "@/components/modals/Command/common/CommandCommandField.tsx";
import { CommandComputedPath } from "@/components/modals/Command/common/CommandComputedPath.tsx";
import { CommandNameField } from "@/components/modals/Command/common/CommandNameField.tsx";
import { CommandWorkingDirectoryField } from "@/components/modals/Command/common/CommandWorkingDirectoryField.tsx";
import {
  formSchema,
  type FormSchemaType,
} from "@/components/modals/Command/common/formSchema.ts";
import { Button } from "@/components/ui/button.tsx";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog.tsx";
import { Form } from "@/components/ui/form.tsx";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { useProjectStore } from "@/store/projectStore.ts";
import { createCommand } from "@/useCases/command/createCommand.ts";

export const CreateCommandModal = ({
  open,
  setOpen,
}: {
  open: boolean;
  setOpen: (open: boolean) => void;
}) => {
  const projectId = useProjectStore((state) => state.projectInfo?.id);

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      command: "",
      workingDirectory: "",
    },
  });

  const onSubmit = async (values: FormSchemaType) => {
    if (!projectId) {
      return;
    }

    try {
      await createCommand({
        id: crypto.randomUUID(),
        projectId: projectId,
        name: values.name,
        command: values.command,
        workingDirectory: values.workingDirectory,
        position: 0, // Will be set by the backend
      });
      toast.success("Command created successfully");

      setOpen(false);
      form.reset();
    } catch (e) {
      toast.error("Failed to create command: " + parseError(e));
    } finally {
      fetchCommands();
    }
  };

  const onOpenChange = (open: boolean) => {
    setOpen(open);
    if (!open) {
      form.reset();
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[628px]">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="w-full">
            <DialogHeader className="flex flex-row items-center gap-6">
              <Terminal />
              <DialogTitle>Create new command</DialogTitle>
            </DialogHeader>
            <div className="space-y-6 my-4">
              <CommandNameField />
              <CommandCommandField />
              <CommandWorkingDirectoryField />
              <CommandComputedPath />
            </div>
            <DialogFooter>
              <DialogClose asChild>
                <Button type="button" variant="outline">
                  Cancel
                </Button>
              </DialogClose>
              <Button type="submit">Create</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};
