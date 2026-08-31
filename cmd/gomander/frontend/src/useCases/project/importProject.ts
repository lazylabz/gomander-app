import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import type { ProjectBlueprint } from "@/contracts/types.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchAvailableProjects } from "@/queries/fetchAvailableProjects.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";

export const importProject = async (
	project: ProjectBlueprint,
	name: string,
	workingDirectory: string,
): Promise<boolean> => {
	try {
		await dataService.importProject(project, name, workingDirectory);

		toast.success(i18n.t("toast.project.importSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.project.importFailed")));
		return false;
	} finally {
		await refreshAfterMutation(fetchAvailableProjects);
	}
};
