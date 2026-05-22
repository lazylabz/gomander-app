import {
	createContext,
	useCallback,
	useContext,
	useEffect,
	useState,
} from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { helpersService } from "@/contracts/service.ts";
import { parseError } from "@/helpers/errorHelpers.ts";

export type UpdateStatus = "idle" | "downloading" | "downloaded" | "installing";

type VersionContext = {
	currentVersion: string;
	newVersion: string | null;
	errorLoadingNewVersion: Error | null;
	updateStatus: UpdateStatus;
	downloadLatestRelease: () => Promise<void>;
	installLatestRelease: () => Promise<void>;
};

export const versionContext = createContext<VersionContext>({
	currentVersion: "",
	newVersion: null,
	errorLoadingNewVersion: null,
	updateStatus: "idle",
	downloadLatestRelease: async () => {},
	installLatestRelease: async () => {},
});

export const VersionProvider = ({
	children,
}: {
	children: React.ReactNode;
}) => {
	const { t } = useTranslation();
	const [currentRelease, setCurrentRelease] = useState<string>("");
	const [newRelease, setNewRelease] = useState<string | null>(null);
	const [errorLoadingNewVersion, setErrorLoadingNewVersion] =
		useState<Error | null>(null);
	const [updateStatus, setUpdateStatus] = useState<UpdateStatus>("idle");
	const [downloadedBinaryPath, setDownloadedBinaryPath] = useState<
		string | null
	>(null);

	const fetchCurrentRelease = useCallback(async () => {
		const release = await helpersService.getCurrentRelease();
		setCurrentRelease(release);
	}, []);

	const checkNewRelease = useCallback(async () => {
		setErrorLoadingNewVersion(null);
		try {
			const release = await helpersService.isThereANewRelease();
			if (release) {
				setNewRelease(release);
			}
		} catch (err) {
			setErrorLoadingNewVersion(err as Error);
			console.error("Error checking for new releases:", err);
			toast.error(t("toast.version.checkError"));
		}
	}, [t]);

	useEffect(() => {
		fetchCurrentRelease();
		checkNewRelease();
	}, [checkNewRelease, fetchCurrentRelease]);

	const downloadLatestRelease = useCallback(async () => {
		if (!newRelease) return;

		setUpdateStatus("downloading");
		try {
			const binaryPath = await helpersService.downloadLatestRelease(newRelease);
			setDownloadedBinaryPath(binaryPath);
			setUpdateStatus("downloaded");
		} catch (err) {
			setUpdateStatus("idle");
			toast.error(parseError(err, t("toast.version.downloadFailed")));
		}
	}, [newRelease, t]);

	const installLatestRelease = useCallback(async () => {
		if (!downloadedBinaryPath) return;

		setUpdateStatus("installing");
		try {
			await helpersService.installLatestReleaseAndQuit(downloadedBinaryPath);
		} catch (err) {
			setUpdateStatus("downloaded");
			toast.error(parseError(err, t("toast.version.installFailed")));
		}
	}, [downloadedBinaryPath, t]);

	return (
		<versionContext.Provider
			value={{
				currentVersion: currentRelease,
				newVersion: newRelease,
				errorLoadingNewVersion,
				updateStatus,
				downloadLatestRelease,
				installLatestRelease,
			}}
		>
			{children}
		</versionContext.Provider>
	);
};

export const useVersionContext = () => useContext(versionContext);
