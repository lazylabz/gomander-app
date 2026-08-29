import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { refreshAfterMutation } from "@/queries/refreshAfterMutation.ts";
import { commandStore } from "@/store/commandStore.ts";

export const reorderCommands = async (
	orderedCommandIds: string[],
): Promise<boolean> => {
	const { commands, setCommands } = commandStore.getState();

	// Optimistic approach to avoid flickering
	setCommands(
		[...commands].sort(
			(a, b) =>
				orderedCommandIds.indexOf(a.id) - orderedCommandIds.indexOf(b.id),
		),
	);

	try {
		await dataService.reorderCommands(orderedCommandIds);
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.command.reorderFailed")));
		return false;
	} finally {
		await refreshAfterMutation(fetchCommands);
	}
};
