import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { releaseStore } from "@/store/releaseStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { resetStores } from "@/testing/stores.ts";
import { checkForNewRelease } from "@/useCases/release/checkForNewRelease.ts";

describe("checkForNewRelease", () => {
	const sut = checkForNewRelease;

	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should store the new release", async () => {
		// Arrange
		installInMemoryBackend({ newRelease: "v2.0.0" });

		// Act
		await sut();

		// Assert
		expect(releaseStore.getState().newVersion).toBe("v2.0.0");
	});

	it("Should store no new version when the backend reports none", async () => {
		// Arrange
		installInMemoryBackend({ newRelease: "" });
		releaseStore.setState({ newVersion: "v2.0.0" });

		// Act
		await sut();

		// Assert
		expect(releaseStore.getState().newVersion).toBeNull();
	});

	it("Should flag the failure and notify the user when the check is rejected", async () => {
		// Arrange
		const consoleError = vi
			.spyOn(console, "error")
			.mockImplementation(() => {});
		const backend = installInMemoryBackend();
		backend.data.checkForNewRelease = async () => {
			throw new Error("boom");
		};

		// Act
		await sut();

		// Assert
		expect(releaseStore.getState().checkFailed).toBe(true);
		expect(toastError).toHaveBeenCalledWith("toast.version.checkError");

		consoleError.mockRestore();
	});

	it("Should clear the failure when a later check succeeds", async () => {
		// Arrange
		installInMemoryBackend({ newRelease: "v2.0.0" });
		releaseStore.setState({ checkFailed: true });

		// Act
		await sut();

		// Assert
		expect(releaseStore.getState().checkFailed).toBe(false);
		expect(toastError).not.toHaveBeenCalled();
	});
});
