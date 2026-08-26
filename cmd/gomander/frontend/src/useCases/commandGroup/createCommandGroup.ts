import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";
import {
	type CommandGroupWithCommandIds,
	resolveGroupCommands,
} from "@/useCases/commandGroup/resolveGroupCommands.ts";

export const createCommandGroup = async (
	args: Omit<CommandGroupWithCommandIds, "position">,
): Promise<boolean> => {
	try {
		await dataService.createCommandGroup({
			id: args.id,
			projectId: args.projectId,
			name: args.name,
			commands: resolveGroupCommands(args.commands),
			position: 0, // Will be set by the backend
		});

		toast.success(i18n.t("toast.commandGroup.createSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.commandGroup.createFailed")));
		return false;
	} finally {
		await refreshAfterMutation(fetchCommandGroups);
	}
};
