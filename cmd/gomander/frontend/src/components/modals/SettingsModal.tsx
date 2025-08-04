import { zodResolver } from "@hookform/resolvers/zod";
import { Plus, Settings, Trash } from "lucide-react";
import { useFieldArray, useForm } from "react-hook-form";
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
import { useUserConfigurationStore } from "@/store/userConfigurationStore.ts";
import { saveUserConfig } from "@/useCases/userConfig/saveUserConfig.ts";

const formSchema = z.object({
  environmentPaths: z.array(
    z.object({
      value: z.string().min(1, "Path cannot be empty"),
    }),
  ),
});

type FormType = z.infer<typeof formSchema>;

export const SettingsModal = ({
  open,
  setOpen,
}: {
  open: boolean;
  setOpen: (open: boolean) => void;
}) => {
  const userConfig = useUserConfigurationStore((state) => state.userConfig);

  const form = useForm<FormType>({
    resolver: zodResolver(formSchema),
    values: {
      environmentPaths:
        userConfig.environmentPaths.map((p) => ({ value: p })) || [],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "environmentPaths" as const,
  });

  const addNewPath = () => {
    append({ value: "" });
  };

  const removePath = (index: number) => () => {
    remove(index);
  };

  const onSubmit = async (data: FormType) => {
    await saveUserConfig({
      lastOpenedProjectId: userConfig.lastOpenedProjectId,
      environmentPaths: data.environmentPaths.map((path) => path.value),
    });

    form.reset();

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
      <DialogContent className="sm:max-w-[628px]">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="w-full">
            <DialogHeader className="flex flex-row items-center gap-4">
              <Settings />
              <DialogTitle>Settings</DialogTitle>
            </DialogHeader>
            <div className="my-4">
              <FormLabel className="mt-8 mb-4">Environment paths</FormLabel>
              {fields.length === 0 && (
                <p className="text-sm text-muted-foreground">
                  Add extra environment paths to your system PATH. These paths
                  will be used to resolve commands and executables.
                </p>
              )}
              {fields.length !== 0 && (
                <div className="flex flex-col gap-2">
                  {fields.map((_, index) => (
                    <div
                      className="flex flex-row w-full items-center gap-3 justify-between"
                      key={index}
                    >
                      <FormField
                        control={form.control}
                        name={`environmentPaths.${index}.value` as const}
                        render={({ field }) => (
                          <FormItem className="flex-1">
                            <FormControl>
                              <Input
                                placeholder="/Users/hackerman/.nvm/versions/node/v20.18.0/bin"
                                {...field}
                              />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <Trash
                        size={18}
                        className="text-muted-foreground cursor-pointer hover:text-destructive"
                        onClick={removePath(index)}
                      />
                    </div>
                  ))}
                </div>
              )}
              <Button
                size="sm"
                type="button"
                variant="ghost"
                className="mt-2"
                onClick={addNewPath}
              >
                Add
                <Plus />
              </Button>
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
