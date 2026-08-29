import { toast } from "sonner";

import { clipboardService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";

export const copyTextToClipboard = async (text: string): Promise<boolean> => {
	const failureMessage = i18n.t("toast.clipboard.copyFailed");

	try {
		const copied = await clipboardService.setText(text);
		if (copied) {
			return true;
		}

		console.error("Failed to copy text to the clipboard");
		toast.error(failureMessage);
		return false;
	} catch (error) {
		console.error("Failed to copy text to the clipboard", error);
		toast.error(parseError(error, failureMessage));
		return false;
	}
};
