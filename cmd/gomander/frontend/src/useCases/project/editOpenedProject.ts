import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import type { Project } from "@/contracts/types.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchAvailableProjects } from "@/queries/fetchAvailableProjects.ts";
import { fetchProject } from "@/queries/fetchProject.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";

export const editOpenedProject = async (project: Project): Promise<boolean> => {
	try {
		await dataService.editProject(project);

		toast.success(i18n.t("toast.settings.projectSaveSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.settings.projectSaveFailed")));
		return false;
	} finally {
		// A rename shows up in the project selection list as well.
		await refreshAfterMutation(fetchProject, fetchAvailableProjects);
	}
};
