import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import type { Command } from "@/contracts/types.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";

export const createCommand = async (command: Command): Promise<boolean> => {
	try {
		await dataService.addCommand(command);

		toast.success(i18n.t("toast.command.createSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.command.createFailed")));
		return false;
	} finally {
		await refreshAfterMutation(fetchCommands);
	}
};
