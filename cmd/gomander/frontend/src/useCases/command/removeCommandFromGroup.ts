import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";

export const removeCommandFromGroup = async (
	commandId: string,
	groupId: string,
): Promise<boolean> => {
	const { commandGroups } = commandGroupStore.getState();
	const group = commandGroups.find((g) => g.id === groupId);

	if (!group) {
		toast.error(i18n.t("toast.commandGroup.notFound"));
		return false;
	}
	if (group.commands.length === 1) {
		toast.error(i18n.t("toast.commandGroup.cannotRemoveLast"));
		return false;
	}

	try {
		await dataService.removeCommandFromGroup(commandId, groupId);

		toast.success(i18n.t("toast.command.removeFromGroupSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.command.removeFromGroupFailed")));
		return false;
	} finally {
		await refreshAfterMutation(fetchCommandGroups);
	}
};
