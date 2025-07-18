import { zodResolver } from "@hookform/resolvers/zod";
import { Terminal } from "lucide-react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button.tsx";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog.tsx";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { useDataContext } from "@/contexts/DataContext.tsx";
import type { Command } from "@/types/contracts.ts";

const formSchema = z.object({
  name: z.string().min(1, {
    message: "Command name is required",
  }),
  command: z.string().min(1, {
    message: "Command is required",
  }),
});

export const EditCommandModal = ({
  command,
  open,
  setOpen,
}: {
  command: Command | null;
  open: boolean;
  setOpen: (open: boolean) => void;
}) => {
  const { editCommand } = useDataContext();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    values: {
      name: command?.name || "",
      command: command?.command || "",
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    if (!command) {
      return;
    }

    await editCommand({
      id: command.id,
      name: values.name,
      command: values.command,
    });

    setOpen(false);
  };

  const onOpenChange = (open: boolean) => {
    setOpen(open);
    if (!open) {
      form.reset();
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="w-full">
            <DialogHeader className="flex flex-row items-center gap-6">
              <Terminal />
              <DialogTitle>Edit command</DialogTitle>
            </DialogHeader>
            <div className="space-y-6 my-4">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder="My awesome command" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="command"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Command</FormLabel>
                    <FormControl>
                      <Input placeholder={'cowsay "Hello World!"'} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
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
