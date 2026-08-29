import { toast } from "sonner";

import { ProjectExportSuccessToast } from "@/components/ProjectExport/ProjectExportSuccessToast.tsx";
import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";

export const exportProject = async (projectId: string): Promise<boolean> => {
	try {
		const exportFilePath = await dataService.exportProject(projectId);

		// The toast offers to open the containing folder, and dismisses itself once
		// it has, so it needs to know its own id.
		const toastId = crypto.randomUUID();
		toast.success(
			<ProjectExportSuccessToast
				exportFilePath={exportFilePath}
				toastId={toastId}
			/>,
			{ id: toastId },
		);
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.project.exportFailed")));
		return false;
	}
};
