import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import type { ProjectBlueprint } from "@/contracts/types.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";

export type ImportSource = "exportedProject" | "packageJson";

export const pickProjectToImport = async (
	source: ImportSource,
): Promise<ProjectBlueprint | null> => {
	try {
		return source === "packageJson"
			? await dataService.getProjectToImportFromPackageJson()
			: await dataService.getProjectToImport();
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.project.selectFailed")));
		return null;
	}
};
