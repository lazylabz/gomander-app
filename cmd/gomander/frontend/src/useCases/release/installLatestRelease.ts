import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { releaseStore } from "@/store/releaseStore.ts";

export const installLatestRelease = async (): Promise<void> => {
	const { downloadedBinaryPath } = releaseStore.getState();
	if (!downloadedBinaryPath) {
		return;
	}

	releaseStore.setState({ updateStatus: "installing" });
	try {
		await dataService.installReleaseAndQuit(downloadedBinaryPath);
	} catch (e) {
		releaseStore.setState({ updateStatus: "downloaded" });
		toast.error(parseError(e, i18n.t("toast.version.installFailed")));
	}
};
