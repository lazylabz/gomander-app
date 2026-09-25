import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { fetchCurrentRelease } from "@/queries/fetchCurrentRelease.ts";
import { releaseStore } from "@/store/releaseStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { resetStores } from "@/testing/stores.ts";

describe("fetchCurrentRelease", () => {
	const sut = fetchCurrentRelease;

	beforeEach(() => {
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should store the current release", async () => {
		// Arrange
		installInMemoryBackend({ currentRelease: "v1.2.3" });

		// Act
		await sut();

		// Assert
		expect(releaseStore.getState().currentVersion).toBe("v1.2.3");
	});
});
