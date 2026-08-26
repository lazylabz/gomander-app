import type { ProjectSettingsSchemaType } from "@/screens/SettingsScreen/schemas/projectSettingsSchema.ts";
import { projectStore } from "@/store/projectStore.ts";
import { editOpenedProject } from "@/useCases/project/editOpenedProject.ts";

export const saveProjectSettingsForm = async (
	formData: ProjectSettingsSchemaType,
) => {
	const { projectInfo } = projectStore.getState();
	if (!projectInfo) {
		return;
	}

	await editOpenedProject({
		...projectInfo,
		name: formData.name,
		workingDirectory: formData.baseWorkingDirectory,
	});
};
