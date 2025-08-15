import { zodResolver } from "@hookform/resolvers/zod";
import { Group } from "lucide-react";
import { useForm } from "react-hook-form";

import { CommandGroupCommandsField } from "@/components/modals/CommandGroup/common/CommandGroupCommandsField.tsx";
import { CommandGroupNameField } from "@/components/modals/CommandGroup/common/CommandGroupNameField.tsx";
import {
  formSchema,
  type FormSchemaType,
} from "@/components/modals/CommandGroup/common/formSchema.ts";
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
import type { CommandGroup } from "@/contracts/types.ts";

export const EditCommandGroupModal = ({
  commandGroup,
  open,
  setOpen,
}: {
  commandGroup: CommandGroup | null;
  open: boolean;
  setOpen: (open: boolean) => void;
}) => {
  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    values: {
      name: commandGroup?.name || "",
      commands: commandGroup?.commands.map((c) => c.id) || [],
    },
  });

  const onSubmit = async (/*values: FormSchemaType*/) => {
    // if (!commandGroup) {
    //   return;
    // }
    //
    // const editedCommandGroup = {
    //   ...commandGroup,
    //   name: values.name,
    //   commands: values.commands,
    // };
    //
    // await editCommandGroup(editedCommandGroup);

    setOpen(false);
    form.reset();
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
            <DialogHeader className="flex flex-row items-center gap-2">
              <Group />
              <DialogTitle>Edit command group</DialogTitle>
            </DialogHeader>
            <div className="space-y-6 my-4">
              <CommandGroupNameField />
              <CommandGroupCommandsField />
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
