import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";

export const stopCommandGroup = async (groupId: string): Promise<boolean> => {
	try {
		await dataService.stopCommandGroup(groupId);
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.commandGroup.stopFailed")));
		return false;
	}
};
