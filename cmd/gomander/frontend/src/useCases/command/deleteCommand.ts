import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";
import { clearLogTail } from "@/store/commandLogsTail.ts";
import { commandStore } from "@/store/commandStore.ts";
import { terminalStore } from "@/store/terminalStore.ts";

export const deleteCommand = async (commandId: string): Promise<boolean> => {
	try {
		await dataService.removeCommand(commandId);

		terminalStore.getState().dispose(commandId);
		clearLogTail(commandId);
		commandStore.getState().setActiveCommandId(null);

		toast.success(i18n.t("toast.command.deleteSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.command.deleteFailed")));
		return false;
	} finally {
		// The backend prunes the command from every group that held it.
		await refreshAfterMutation(fetchCommands, fetchCommandGroups);
	}
};
