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
import type { Command } from "@/contracts/types.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { editCommand } from "@/useCases/command/editCommand.ts";

export const EditCommandModal = ({
  command,
  open,
  setOpen,
}: {
  command: Command | null;
  open: boolean;
  setOpen: (open: boolean) => void;
}) => {
  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    values: {
      name: command?.name || "",
      command: command?.command || "",
      workingDirectory: command?.workingDirectory || "",
    },
  });

  const onSubmit = async (values: FormSchemaType) => {
    if (!command) {
      return;
    }

    try {
      await editCommand({
        ...command,
        // Editable fields
        name: values.name,
        command: values.command,
        workingDirectory: values.workingDirectory,
      });

      toast.success("Command updated successfully");

      setOpen(false);
    } catch (e: unknown) {
      toast.error("Failed to update command: " + parseError(e));
    } finally {
      fetchCommands();
      fetchCommandGroups();
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
              <DialogTitle>Edit command</DialogTitle>
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
              <Button type="submit">Save</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};
