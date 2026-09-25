import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { releaseStore } from "@/store/releaseStore.ts";

export const checkForNewRelease = async (): Promise<void> => {
	releaseStore.setState({ checkFailed: false });
	try {
		const release = await dataService.checkForNewRelease();
		releaseStore.setState({ newVersion: release || null });
	} catch (e) {
		console.error("Error checking for new releases:", e);
		releaseStore.setState({ checkFailed: true });
		toast.error(i18n.t("toast.version.checkError"));
	}
};
