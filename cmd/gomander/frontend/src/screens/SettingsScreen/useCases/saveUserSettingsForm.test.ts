import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import type { Localization } from "@/contracts/types.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { saveUserSettingsForm } from "@/screens/SettingsScreen/useCases/saveUserSettingsForm.ts";
import { userConfigurationStore } from "@/store/userConfigurationStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { UserConfigBuilder } from "@/testing/builders/userConfig.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { resetStores } from "@/testing/stores.ts";

describe("saveUserSettingsForm", () => {
	const sut = saveUserSettingsForm;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	const environmentPath = { id: crypto.randomUUID(), path: "/usr/local/bin" };
	const userConfig = new UserConfigBuilder()
		.withLastOpenedProjectId("project-1")
		.withLocale("en")
		.build();
	const spanish = {
		toast: { settings: { userSaveSuccess: "Ajustes guardados" } },
	} as unknown as Localization;

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		resetStores();
		userConfigurationStore.setState({ userConfig });
	});

	afterEach(async () => {
		resetBackendServices();
		// i18n is a module singleton: a language or bundle left behind would make
		// the next test translate keys instead of echoing them.
		await i18n.changeLanguage("en");
		i18n.removeResourceBundle("es", "translation");
	});

	it("Should save the settings keeping the last opened project and refresh the user configuration", async () => {
		// Arrange
		const backend = installInMemoryBackend({ userConfig });
		const saved = {
			lastOpenedProjectId: "project-1",
			environmentPaths: [environmentPath],
			locale: "en",
		};

		// Act
		await sut({ environmentPaths: [environmentPath], locale: "en" });

		// Assert
		expect(backend.state.userConfig).toEqual(saved);
		expect(userConfigurationStore.getState().userConfig).toEqual(saved);
		expect(toastSuccess).toHaveBeenCalledWith("toast.settings.userSaveSuccess");
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should switch the interface to the chosen language loading its translations from the backend", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			userConfig,
			translations: { es: spanish },
		});

		// Act
		await sut({ environmentPaths: [], locale: "es" });

		// Assert
		expect(i18n.language).toBe("es");
		expect(backend.state.userConfig.locale).toBe("es");
		expect(toastSuccess).toHaveBeenCalledWith("Ajustes guardados");
	});

	it("Should notify the user and still refresh the user configuration when the backend rejects the save", async () => {
		// Arrange
		const stored = new UserConfigBuilder()
			.withLastOpenedProjectId("project-1")
			.withEnvironmentPaths({ id: crypto.randomUUID(), path: "/stored" })
			.build();
		const backend = installInMemoryBackend({ userConfig: stored });
		backend.data.saveUserConfig = async () => {
			throw new Error("boom");
		};

		// Act
		await sut({ environmentPaths: [environmentPath], locale: "en" });

		// Assert
		expect(toastError).toHaveBeenCalledWith(
			"toast.settings.userSaveFailed: boom",
		);
		expect(toastSuccess).not.toHaveBeenCalled();
		expect(userConfigurationStore.getState().userConfig).toEqual(stored);
	});
});
