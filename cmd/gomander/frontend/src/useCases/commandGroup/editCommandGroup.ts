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

export const editCommandGroup = async (
	args: CommandGroupWithCommandIds,
): Promise<boolean> => {
	try {
		await dataService.editCommandGroup({
			id: args.id,
			projectId: args.projectId,
			name: args.name,
			commands: resolveGroupCommands(args.commands),
			position: args.position,
		});

		toast.success(i18n.t("toast.commandGroup.updateSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.commandGroup.updateFailed")));
		return false;
	} finally {
		await refreshAfterMutation(fetchCommandGroups);
	}
};
