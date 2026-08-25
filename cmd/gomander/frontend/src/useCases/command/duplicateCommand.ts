import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";
import { commandStore } from "@/store/commandStore.ts";

export const duplicateCommand = async (
	commandId: string,
	targetGroupId?: string,
): Promise<boolean> => {
	try {
		await dataService.duplicateCommand(commandId, targetGroupId || "");

		commandStore.getState().setActiveCommandId(null);

		toast.success(i18n.t("toast.command.duplicateSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.command.duplicateFailed")));
		return false;
	} finally {
		await refreshAfterMutation(fetchCommands, fetchCommandGroups);
	}
};
