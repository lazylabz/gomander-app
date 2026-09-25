import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { releaseStore } from "@/store/releaseStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { resetStores } from "@/testing/stores.ts";
import { installLatestRelease } from "@/useCases/release/installLatestRelease.ts";

describe("installLatestRelease", () => {
	const sut = installLatestRelease;

	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should install nothing when no release was downloaded", async () => {
		// Arrange
		const backend = installInMemoryBackend();

		// Act
		await sut();

		// Assert
		expect(backend.state.installedBinaryPath).toBeNull();
		expect(releaseStore.getState().updateStatus).toBe("idle");
	});

	it("Should install the downloaded binary", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		releaseStore.setState({
			updateStatus: "downloaded",
			downloadedBinaryPath: "/downloads/gomander",
		});

		// Act
		await sut();

		// Assert
		expect(backend.state.installedBinaryPath).toBe("/downloads/gomander");
		expect(releaseStore.getState().updateStatus).toBe("installing");
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should go back to downloaded and notify the user when the install is rejected", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.installReleaseAndQuit = async () => {
			throw new Error("boom");
		};
		releaseStore.setState({
			updateStatus: "downloaded",
			downloadedBinaryPath: "/downloads/gomander",
		});

		// Act
		await sut();

		// Assert
		expect(releaseStore.getState().updateStatus).toBe("downloaded");
		expect(releaseStore.getState().downloadedBinaryPath).toBe(
			"/downloads/gomander",
		);
		expect(toastError).toHaveBeenCalledWith(
			"toast.version.installFailed: boom",
		);
	});
});
