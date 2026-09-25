import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

export type UpdateStatus = "idle" | "downloading" | "downloaded" | "installing";

type ReleaseStore = {
	currentRelease: string;
	newRelease: string | null;
	checkFailed: boolean;
	updateStatus: UpdateStatus;
	downloadedBinaryPath: string | null;
};

export const releaseStore = createStore<ReleaseStore>()(() => ({
	currentRelease: "",
	newRelease: null,
	checkFailed: false,
	updateStatus: "idle",
	downloadedBinaryPath: null,
}));

export const useReleaseStore = <T>(selector: (state: ReleaseStore) => T): T =>
	useStore(releaseStore, selector);
