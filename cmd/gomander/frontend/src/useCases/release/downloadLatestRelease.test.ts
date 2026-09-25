import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { releaseStore } from "@/store/releaseStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { resetStores } from "@/testing/stores.ts";
import { downloadLatestRelease } from "@/useCases/release/downloadLatestRelease.ts";

describe("downloadLatestRelease", () => {
	const sut = downloadLatestRelease;

	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should download nothing when there is no new version", async () => {
		// Arrange
		const backend = installInMemoryBackend();

		// Act
		await sut();

		// Assert
		expect(backend.state.downloadedReleases).toEqual([]);
		expect(releaseStore.getState().updateStatus).toBe("idle");
	});

	it("Should download the new version and keep the path of its binary", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			releaseBinaryPath: "/downloads/gomander",
		});
		releaseStore.setState({ newRelease: "v2.0.0" });

		// Act
		await sut();

		// Assert
		expect(backend.state.downloadedReleases).toEqual(["v2.0.0"]);
		expect(releaseStore.getState().updateStatus).toBe("downloaded");
		expect(releaseStore.getState().downloadedBinaryPath).toBe(
			"/downloads/gomander",
		);
	});

	it("Should be downloading while the download is in flight", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		let finishDownload: (path: string) => void = () => {};
		backend.data.downloadRelease = () =>
			new Promise((resolve) => {
				finishDownload = resolve;
			});
		releaseStore.setState({ newRelease: "v2.0.0" });

		// Act
		const download = sut();

		// Assert
		expect(releaseStore.getState().updateStatus).toBe("downloading");
		finishDownload("/downloads/gomander");
		await download;
	});

	it("Should go back to idle and notify the user when the download is rejected", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.downloadRelease = async () => {
			throw new Error("boom");
		};
		releaseStore.setState({ newRelease: "v2.0.0" });

		// Act
		await sut();

		// Assert
		expect(releaseStore.getState().updateStatus).toBe("idle");
		expect(releaseStore.getState().downloadedBinaryPath).toBeNull();
		expect(toastError).toHaveBeenCalledWith(
			"toast.version.downloadFailed: boom",
		);
	});
});
