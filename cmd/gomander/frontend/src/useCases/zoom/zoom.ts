import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import type { ZoomDirection } from "@/contracts/types.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";

const ZOOM_BY_DIRECTION: Record<ZoomDirection, () => Promise<void>> = {
	in: () => dataService.zoomIn(),
	out: () => dataService.zoomOut(),
	reset: () => dataService.resetZoom(),
};

export const zoom = async (direction: ZoomDirection): Promise<boolean> => {
	try {
		await ZOOM_BY_DIRECTION[direction]();
		return true;
	} catch (e) {
		toast.error(parseError(e, i18n.t("toast.settings.zoomFailed")));
		return false;
	}
};
