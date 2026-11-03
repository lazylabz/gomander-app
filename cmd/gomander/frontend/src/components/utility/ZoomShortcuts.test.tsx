import { act, render } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it } from "vitest";

import { ZoomShortcuts } from "@/components/utility/ZoomShortcuts.tsx";
import { resetBackendServices } from "@/contracts/service.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";

describe("ZoomShortcuts", () => {
	// The listener goes on once the backend has said which OS this is.
	const renderSut = async () => {
		await act(async () => {
			render(<ZoomShortcuts />);
		});
	};

	afterEach(() => {
		resetBackendServices();
	});

	it("Should zoom from Ctrl and +, - or 0 outside macOS", async () => {
		// Arrange
		const backend = installInMemoryBackend({ os: "windows" });
		const user = userEvent.setup();
		await renderSut();

		// Act
		await user.keyboard("{Control>}={/Control}");
		await user.keyboard("{Control>}-{/Control}");
		await user.keyboard("{Control>}0{/Control}");

		// Assert
		expect(backend.state.zoomRequests).toEqual(["in", "out", "reset"]);
	});

	it("Should leave the keys to the native shortcuts on macOS", async () => {
		// Arrange
		const backend = installInMemoryBackend({ os: "darwin" });
		const user = userEvent.setup();
		await renderSut();

		// Act
		await user.keyboard("{Control>}={/Control}");

		// Assert
		expect(backend.state.zoomRequests).toEqual([]);
	});
});
