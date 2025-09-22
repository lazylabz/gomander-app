import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";

import { BaseWorkingDirectoryField } from "@/components/modals/Project/common/BaseWorkingDirectoryField.tsx";
import {
  formSchema,
  type FormSchemaType,
} from "@/components/modals/Project/common/importAndExportSchema.ts";
import { ProjectCommandGroupsField } from "@/components/modals/Project/common/ProjectCommandGroupsField.tsx";
import { ProjectCommandsField } from "@/components/modals/Project/common/ProjectCommandsField.tsx";
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
      commands: project?.commands.map((c) => c.id) || [],
      commandGroups: project?.commandGroups.map((c) => c.id) || [],
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

    const projectWithSelectedCommandsAndCommandGroups: ProjectExport = {
      ...project,
      commands: project.commands.filter((c) => values.commands.includes(c.id)),
      commandGroups: project.commandGroups.filter((cg) =>
        values.commandGroups.includes(cg.id),
      ),
    };

    try {
      await importProject(
        projectWithSelectedCommandsAndCommandGroups,
        values.name,
        values.baseWorkingDirectory,
      );

      await onSuccess();
      handleOpenChange(false);
      toast.success("Project imported successfully");
    } catch (e) {
      toast.error(parseError(e, "Failed to import the project"));
    }
  };

  const commandIdsWatcher = form.watch("commands");

  useEffect(() => {
    if (!project) return;

    const currentCommandGroups = form.getValues("commandGroups");

    const updatedCommandGroups = currentCommandGroups.filter((groupId) => {
      const group = project.commandGroups.find((cg) => cg.id === groupId);
      // Keep the group checked only if at least one of its commands is selected
      return (
        group &&
        group.commandIds.some((commandId) =>
          commandIdsWatcher.includes(commandId),
        )
      );
    });

    // Only update if there's a difference to avoid infinite loops
    if (currentCommandGroups.length !== updatedCommandGroups.length) {
      form.setValue("commandGroups", updatedCommandGroups);
    }
  }, [commandIdsWatcher, project, form]);

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="min-w-[600px]">
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
              <div className="flex items-start flex-wrap flex-col sm:flex-row gap-4 justify-between">
                <ProjectCommandsField commands={project?.commands || []} />
                <ProjectCommandGroupsField
                  commandGroups={project?.commandGroups || []}
                  selectedCommandIds={commandIdsWatcher}
                />
              </div>

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
