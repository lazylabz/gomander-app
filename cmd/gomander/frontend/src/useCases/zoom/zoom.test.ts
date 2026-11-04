import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { zoom } from "@/useCases/zoom/zoom.ts";

describe("zoom", () => {
	const sut = zoom;

	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should ask the backend to zoom in, out and back to 100%", async () => {
		// Arrange
		const backend = installInMemoryBackend();

		// Act
		await sut("in");
		await sut("out");
		await sut("reset");

		// Assert
		expect(backend.state.zoomRequests).toEqual(["in", "out", "reset"]);
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should tell the user when the backend cannot zoom", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.zoomIn = async () => {
			throw new Error("boom");
		};

		// Act
		const zoomed = await sut("in");

		// Assert
		expect(zoomed).toBe(false);
		expect(toastError).toHaveBeenCalledWith("toast.settings.zoomFailed: boom");
	});
});
