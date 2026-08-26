import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { cleanCommandGroupFromErrors } from "@/useCases/command/cleanCommandGroupFromErrors.ts";

export const runCommandGroup = async (groupId: string): Promise<boolean> => {
	cleanCommandGroupFromErrors(groupId);

	try {
		await dataService.runCommandGroup(groupId);
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.commandGroup.runFailed")));
		return false;
	}
};
