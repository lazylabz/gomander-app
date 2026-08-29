import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";

export const deleteCommandGroup = async (
	commandGroupId: string,
): Promise<boolean> => {
	try {
		await dataService.deleteCommandGroup(commandGroupId);

		toast.success(i18n.t("toast.commandGroup.deleteSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.commandGroup.deleteFailed")));
		return false;
	} finally {
		// The commands the group held outlive it, so only the groups change.
		await refreshAfterMutation(fetchCommandGroups);
	}
};
