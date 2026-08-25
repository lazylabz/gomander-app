import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import type { Command } from "@/contracts/types.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";

export const editCommand = async (command: Command): Promise<boolean> => {
	try {
		await dataService.editCommand(command);

		toast.success(i18n.t("toast.command.updateSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.command.updateFailed")));
		return false;
	} finally {
		// The groups carry a copy of every command they hold.
		await refreshAfterMutation(fetchCommands, fetchCommandGroups);
	}
};
