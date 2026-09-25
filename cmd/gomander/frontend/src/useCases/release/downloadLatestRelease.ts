import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { parseError } from "@/helpers/errorHelpers.ts";
import { releaseStore } from "@/store/releaseStore.ts";

export const downloadLatestRelease = async (): Promise<void> => {
	const { newVersion } = releaseStore.getState();
	if (!newVersion) {
		return;
	}

	releaseStore.setState({ updateStatus: "downloading" });
	try {
		const binaryPath = await dataService.downloadRelease(newVersion);
		releaseStore.setState({
			updateStatus: "downloaded",
			downloadedBinaryPath: binaryPath,
		});
	} catch (e) {
		releaseStore.setState({ updateStatus: "idle" });
		toast.error(parseError(e, i18n.t("toast.version.downloadFailed")));
	}
};
