import { zodResolver } from "@hookform/resolvers/zod";
import { Terminal } from "lucide-react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";

import { CommandCommandField } from "@/components/modals/Command/common/CommandCommandField.tsx";
import { CommandComputedPath } from "@/components/modals/Command/common/CommandComputedPath.tsx";
import { CommandErrorPatternsField } from "@/components/modals/Command/common/CommandErrorPatternsField.tsx";
import { CommandLinkField } from "@/components/modals/Command/common/CommandLinkField.tsx";
import { CommandNameField } from "@/components/modals/Command/common/CommandNameField.tsx";
import { CommandWorkingDirectoryField } from "@/components/modals/Command/common/CommandWorkingDirectoryField.tsx";
import {
  formSchema,
  type FormSchemaType,
} from "@/components/modals/Command/common/formSchema.ts";
import type { Command } from "@/contracts/types.ts";
import { Button } from "@/design-system/components/ui/button.tsx";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/design-system/components/ui/dialog.tsx";
import { Form } from "@/design-system/components/ui/form.tsx";
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
      link: command?.link || "",
      errorPatterns: command?.errorPatterns || "",
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
        link: values.link,
        errorPatterns: values.errorPatterns,
      });

      toast.success("Command updated successfully");

      setOpen(false);
    } catch (e: unknown) {
      toast.error(parseError(e, "Failed to update command"));
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
              <CommandLinkField />
              <CommandErrorPatternsField />
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
