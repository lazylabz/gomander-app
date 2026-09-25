import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchAvailableProjects } from "@/queries/fetchAvailableProjects.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";

export const createProject = async (
	name: string,
	workingDirectory: string,
): Promise<boolean> => {
	try {
		await dataService.createProject({
			id: crypto.randomUUID(),
			name,
			workingDirectory,
		});

		toast.success(i18n.t("toast.project.createSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.project.createFailed")));
		return false;
	} finally {
		await refreshAfterMutation(fetchAvailableProjects);
	}
};
