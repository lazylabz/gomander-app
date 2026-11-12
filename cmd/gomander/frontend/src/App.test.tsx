import { cleanup, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import App from "@/App.tsx";
import type { InMemoryBackend } from "@/contracts/adapters/inMemory.ts";
import { resetBackendServices } from "@/contracts/service.ts";
import type { Localization } from "@/contracts/types.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import { resetStores } from "@/testing/stores.ts";

const echoedKeys = [
	"common.loadFailed",
	"common.retry",
	"projectSelection.emptyState",
	"projectSelection.welcome",
];

describe("App", () => {
	let backend: InMemoryBackend;

	const failCommandsLoad = () => {
		const original = backend.data.getCommands;
		backend.data.getCommands = async () => {
			throw new Error("boom");
		};
		return () => {
			backend.data.getCommands = original;
		};
	};

	const failUserConfigLoad = () => {
		const original = backend.data.getUserConfig;
		backend.data.getUserConfig = async () => {
			throw new Error("boom");
		};
		return () => {
			backend.data.getUserConfig = original;
		};
	};

	beforeEach(() => {
		vi.spyOn(console, "error").mockImplementation(() => {});
		vi.spyOn(console, "warn").mockImplementation(() => {});
		resetStores();
		backend = installInMemoryBackend({
			translations: {
				en: Object.fromEntries(
					echoedKeys.map((key) => [key, key]),
				) as unknown as Localization,
			},
		});
	});

	afterEach(() => {
		cleanup();
		resetBackendServices();
		vi.restoreAllMocks();
	});

	it("Should show the project selection once everything has loaded", async () => {
		// Arrange / Act
		renderWithProviders(<App />);

		// Assert
		expect(
			await screen.findByText("projectSelection.emptyState"),
		).toBeInTheDocument();
	});

	it("Should tell the user the startup failed instead of staying blank when the project data fails to load", async () => {
		// Arrange
		failCommandsLoad();

		// Act
		renderWithProviders(<App />);

		// Assert
		expect(await screen.findByText("common.loadFailed")).toBeInTheDocument();
		expect(
			screen.queryByText("projectSelection.emptyState"),
		).not.toBeInTheDocument();
	});

	it("Should start the app once retrying succeeds", async () => {
		// Arrange
		const recoverBackend = failCommandsLoad();
		const { user } = renderWithProviders(<App />);
		await screen.findByText("common.loadFailed");
		recoverBackend();

		// Act
		await user.click(screen.getByRole("button", { name: "common.retry" }));

		// Assert
		expect(
			await screen.findByText("projectSelection.emptyState"),
		).toBeInTheDocument();
		expect(screen.queryByText("common.loadFailed")).not.toBeInTheDocument();
	});

	it("Should offer the retry in English when not even the translations could be loaded", async () => {
		// Arrange
		failUserConfigLoad();

		// Act
		renderWithProviders(<App />);

		// Assert
		expect(
			await screen.findByText("Something went wrong while loading. Try again."),
		).toBeInTheDocument();
		expect(screen.getByRole("button", { name: "Retry" })).toBeInTheDocument();
	});

	it("Should start the app with the real translations once retrying succeeds after the translations failed", async () => {
		// Arrange
		const recoverBackend = failUserConfigLoad();
		const { user } = renderWithProviders(<App />);
		await user.click(await screen.findByRole("button", { name: "Retry" }));
		await screen.findByRole("button", { name: "Retry" });
		recoverBackend();

		// Act
		await user.click(screen.getByRole("button", { name: "Retry" }));

		// Assert
		expect(
			await screen.findByText("projectSelection.emptyState"),
		).toBeInTheDocument();
	});
});
