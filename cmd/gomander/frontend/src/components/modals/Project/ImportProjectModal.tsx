import { zodResolver } from "@hookform/resolvers/zod";
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
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion.tsx";
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
      commandGroups: project?.commandGroups.map((cg) => cg.id) || [],
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

  const handleCommandIdsChange = (selectedCommandIds: string[]) => {
    const currentCommandGroups = form.getValues("commandGroups");

    const updatedCommandGroups = currentCommandGroups.filter((groupId) => {
      const group = project?.commandGroups.find((cg) => cg.id === groupId);
      // Keep the group checked only if at least one of its commands is selected
      return (
        group &&
        group.commandIds.some((commandId) =>
          selectedCommandIds.includes(commandId),
        )
      );
    });

    form.setValue("commandGroups", updatedCommandGroups);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="w-full sm:max-w-[800px]">
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
              <Accordion type="single" collapsible>
                <AccordionItem value="1">
                  <AccordionTrigger>Advanced import</AccordionTrigger>
                  <AccordionContent>
                    <div className="flex sm:items-start flex-wrap flex-col items-stretch sm:flex-row gap-4 justify-between">
                      <ProjectCommandsField
                        onChange={handleCommandIdsChange}
                        commands={project?.commands || []}
                      />
                      <ProjectCommandGroupsField
                        commandGroups={project?.commandGroups || []}
                        selectedCommandIds={commandIdsWatcher}
                        commands={project?.commands || []}
                      />
                    </div>
                  </AccordionContent>
                </AccordionItem>
              </Accordion>
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
