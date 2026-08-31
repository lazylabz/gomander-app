import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import type { CommandGroup } from "@/contracts/types.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";

export const editCommandGroup = async (
	args: CommandGroup,
): Promise<boolean> => {
	try {
		await dataService.editCommandGroup(args);

		toast.success(i18n.t("toast.commandGroup.updateSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.commandGroup.updateFailed")));
		return false;
	} finally {
		await refreshAfterMutation(fetchCommandGroups);
	}
};
