import { toast } from "sonner";

import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchProject } from "@/queries/fetchProject.ts";
import type { ProjectSettingsSchemaType } from "@/screens/SettingsScreen/schemas/projectSettingsSchema.ts";
import { projectStore } from "@/store/projectStore.ts";
import { editOpenedProject } from "@/useCases/project/editOpenedProject.ts";

export const saveProjectSettingsForm = async (formData: ProjectSettingsSchemaType) => {
  const { projectInfo } = projectStore.getState();
  if (!projectInfo) {
    return;
  }

  try {
    await editOpenedProject({
      ...projectInfo,
      name: formData.name,
      workingDirectory: formData.baseWorkingDirectory,
    });
    toast.success("Project settings saved successfully");
  } catch (e) {
    toast.error(parseError(e, "Failed to save project settings"));
  }

  await fetchProject();
};