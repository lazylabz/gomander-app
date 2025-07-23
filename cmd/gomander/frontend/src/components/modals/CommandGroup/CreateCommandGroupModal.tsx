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
import { useDataContext } from "@/contexts/DataContext.tsx";

export const CreateCommandGroupModal = ({
  open,
  setOpen,
}: {
  open: boolean;
  setOpen: (open: boolean) => void;
}) => {
  const { commandGroups, saveCommandGroups } = useDataContext();

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      commands: [],
    },
  });

  const onSubmit = async (values: FormSchemaType) => {
    await saveCommandGroups([
      ...commandGroups,
      {
        id: crypto.randomUUID(),
        name: values.name,
        commands: values.commands,
      },
    ]);

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
              <DialogTitle>Create new command group</DialogTitle>
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
              <Button type="submit">Create</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};
