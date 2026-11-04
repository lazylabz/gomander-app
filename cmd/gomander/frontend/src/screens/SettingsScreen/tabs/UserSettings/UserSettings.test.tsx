import { screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { SettingsContextProvider } from "@/screens/SettingsScreen/context/settingsContext.tsx";
import { UserSettings } from "@/screens/SettingsScreen/tabs/UserSettings/UserSettings.tsx";
import { userConfigurationStore } from "@/store/userConfigurationStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { UserConfigBuilder } from "@/testing/builders/userConfig.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import { resetStores } from "@/testing/stores.ts";

describe("UserSettings", () => {
	const renderSut = () =>
		renderWithProviders(
			<SettingsContextProvider>
				<UserSettings />
			</SettingsContextProvider>,
		);

	beforeEach(async () => {
		await installTranslations();
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should save the logs font size the user types", async () => {
		// Arrange
		const userConfig = new UserConfigBuilder().withLogsFontSize(15).build();
		const backend = installInMemoryBackend({ userConfig });
		userConfigurationStore.getState().setUserConfig(userConfig);
		const { user } = renderSut();
		const fontSize = screen.getByLabelText(
			"userSettingsForm.logsFontSizeLabel",
		);

		// Act
		await user.clear(fontSize);
		await user.type(fontSize, "20");

		// Assert
		await waitFor(() => expect(backend.state.userConfig.logsFontSize).toBe(20));
	});

	it("Should refuse a logs font size out of range without saving it", async () => {
		// Arrange
		const userConfig = new UserConfigBuilder().withLogsFontSize(15).build();
		const backend = installInMemoryBackend({ userConfig });
		userConfigurationStore.getState().setUserConfig(userConfig);
		const { user } = renderSut();
		const fontSize = screen.getByLabelText(
			"userSettingsForm.logsFontSizeLabel",
		);

		// Act
		await user.clear(fontSize);
		await user.type(fontSize, "99");

		// Assert
		expect(
			await screen.findByText("userSettingsForm.validation.logsFontSize"),
		).toBeInTheDocument();
		// Past the autosave debounce, so a save it let through would have landed
		await new Promise((resolve) => setTimeout(resolve, 400));
		expect(backend.state.userConfig.logsFontSize).toBe(15);
	});
});
