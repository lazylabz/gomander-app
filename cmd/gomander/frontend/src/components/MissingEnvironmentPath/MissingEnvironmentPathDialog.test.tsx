import { screen } from "@testing-library/react";
import { beforeEach, describe, expect, it } from "vitest";

import { MissingEnvironmentPathDialog } from "@/components/MissingEnvironmentPath/MissingEnvironmentPathDialog.tsx";
import { ScreenRoutes } from "@/routes.ts";
import { missingEnvironmentPathStore } from "@/store/missingEnvironmentPathStore.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import { resetStores } from "@/testing/stores.ts";

describe("MissingEnvironmentPathDialog", () => {
	beforeEach(async () => {
		await installTranslations();
		resetStores();
	});

	it("Should stay hidden while the dialog is closed", () => {
		// Arrange / Act
		renderWithProviders(<MissingEnvironmentPathDialog />);

		// Assert
		expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
	});

	it("Should show the dialog when it is opened", () => {
		// Arrange
		missingEnvironmentPathStore.getState().setDialogOpen(true);

		// Act
		renderWithProviders(<MissingEnvironmentPathDialog />);

		// Assert
		expect(
			screen.getByRole("dialog", { name: "modal.missingPath.title" }),
		).toBeInTheDocument();
	});

	it("Should close the dialog and take the user to the user settings", async () => {
		// Arrange
		missingEnvironmentPathStore.getState().setDialogOpen(true);
		const { user, location } = renderWithProviders(
			<MissingEnvironmentPathDialog />,
		);

		// Act
		await user.click(
			screen.getByRole("button", { name: "modal.missingPath.goToSettings" }),
		);

		// Assert
		expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
		expect(missingEnvironmentPathStore.getState().dialogOpen).toBe(false);
		expect(location().pathname).toBe(ScreenRoutes.Settings);
		expect(location().state).toEqual({ tab: "user" });
	});
});
