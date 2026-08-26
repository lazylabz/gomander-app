import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";

export const reorderCommandGroups = async (
	orderedCommandGroupIds: string[],
): Promise<boolean> => {
	const { commandGroups, setCommandGroups } = commandGroupStore.getState();

	// Optimistic approach to avoid flickering
	setCommandGroups(
		[...commandGroups].sort(
			(a, b) =>
				orderedCommandGroupIds.indexOf(a.id) -
				orderedCommandGroupIds.indexOf(b.id),
		),
	);

	try {
		await dataService.reorderCommandGroups(orderedCommandGroupIds);

		toast.success(i18n.t("toast.commandGroup.reorderSuccess"));
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.commandGroup.reorderFailed")));
		return false;
	} finally {
		await refreshAfterMutation(fetchCommandGroups);
	}
};
