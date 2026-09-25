import { afterEach, describe, expect, it } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { fetchSupportedLanguages } from "@/queries/fetchSupportedLanguages.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";

describe("fetchSupportedLanguages", () => {
	const sut = fetchSupportedLanguages;

	afterEach(() => {
		resetBackendServices();
	});

	it("Should label each known language in its own language", async () => {
		// Arrange
		installInMemoryBackend({ supportedLanguages: ["en", "es"] });

		// Act
		const languages = await sut();

		// Assert
		expect(languages).toEqual([
			{ value: "en", label: "English" },
			{ value: "es", label: "Español" },
		]);
	});

	it("Should label an unknown language with its code", async () => {
		// Arrange
		installInMemoryBackend({ supportedLanguages: ["fr"] });

		// Act
		const languages = await sut();

		// Assert
		expect(languages).toEqual([{ value: "fr", label: "fr" }]);
	});
});
