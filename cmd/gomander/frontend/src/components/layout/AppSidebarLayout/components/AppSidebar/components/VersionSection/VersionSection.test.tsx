import { screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { VersionSection } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/VersionSection/VersionSection.tsx";
import { releaseStore } from "@/store/releaseStore.ts";
import { installTranslations, withTranslation } from "@/testing/i18n.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import { resetStores } from "@/testing/stores.ts";

describe("VersionSection", () => {
	const openAboutModal = vi.fn();

	const render = () =>
		renderWithProviders(<VersionSection openAboutModal={openAboutModal} />);

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
		resetStores();
	});

	it("Should show the current version", () => {
		// Arrange
		withTranslation("sidebar.version.current", "{{version}}");
		releaseStore.setState({ currentRelease: "v1.2.0" });

		// Act
		render();

		// Assert
		expect(screen.getByText("v1.2.0")).toBeInTheDocument();
	});

	it("Should show a placeholder while the current version is unknown", () => {
		// Act
		render();

		// Assert
		expect(screen.getByText("...")).toBeInTheDocument();
	});

	it.each([
		{
			situation: "a new version is available",
			release: { currentRelease: "v1.2.0", newRelease: "v1.3.0" },
			tooltip: "sidebar.version.newAvailable",
		},
		{
			situation: "the current version is the latest",
			release: { currentRelease: "v1.2.0" },
			tooltip: "sidebar.version.latest",
		},
		{
			situation: "the check for a new version failed",
			release: { currentRelease: "v1.2.0", checkFailed: true },
			tooltip: "sidebar.version.checkError",
		},
	])("Should tell on hover when $situation", async ({ release, tooltip }) => {
		// Arrange
		releaseStore.setState(release);
		const { user } = render();

		// Act
		await user.hover(screen.getByText("sidebar.version.current"));

		// Assert
		expect(await screen.findByRole("tooltip")).toHaveTextContent(tooltip);
	});

	it("Should open the about modal when clicked", async () => {
		// Arrange
		releaseStore.setState({ currentRelease: "v1.2.0" });
		const { user } = render();

		// Act
		await user.click(screen.getByText("sidebar.version.current"));

		// Assert
		expect(openAboutModal).toHaveBeenCalled();
	});
});
