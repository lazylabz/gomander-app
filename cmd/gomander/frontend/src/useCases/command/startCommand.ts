import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { cleanCommandError } from "@/useCases/command/cleanCommandError.ts";

export const startCommand = async (commandId: string): Promise<boolean> => {
	cleanCommandError(commandId);

	try {
		await dataService.runCommand(commandId);
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.command.runFailed")));
		return false;
	}
};
