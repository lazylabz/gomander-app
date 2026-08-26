import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchAvailableProjects } from "@/queries/fetchAvailableProjects.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";

export const deleteProject = async (projectId: string): Promise<boolean> => {
	try {
		await dataService.deleteProject(projectId);

		toast.success(i18n.t("toast.project.deleteSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.project.deleteFailed")));
		return false;
	} finally {
		await refreshAfterMutation(fetchAvailableProjects);
	}
};
