import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { projectStore } from "@/store/projectStore.ts";

export const closeProject = async (): Promise<boolean> => {
	try {
		await dataService.closeProject();

		projectStore.getState().setProjectInfo(null);
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.project.closeFailed")));
		return false;
	}
};
