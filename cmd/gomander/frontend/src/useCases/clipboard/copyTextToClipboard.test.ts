import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { copyTextToClipboard } from "@/useCases/clipboard/copyTextToClipboard.ts";

describe("copyTextToClipboard", () => {
	const toastError = vi.spyOn(toast, "error");
	const consoleError = vi.spyOn(console, "error").mockImplementation(() => {});

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should copy text through the backend clipboard", async () => {
		// Arrange
		const backend = installInMemoryBackend();

		// Act
		const copied = await copyTextToClipboard("selected output");

		// Assert
		expect(copied).toBe(true);
		expect(backend.state.clipboardText).toBe("selected output");
		expect(toastError).not.toHaveBeenCalled();
	});

	it("Should notify the user when the clipboard rejects the write", async () => {
		// Arrange
		installInMemoryBackend({ clipboardWriteResult: false });

		// Act
		const copied = await copyTextToClipboard("selected output");

		// Assert
		expect(copied).toBe(false);
		expect(consoleError).toHaveBeenCalledWith(
			"Failed to copy text to the clipboard",
		);
		expect(toastError).toHaveBeenCalledWith("toast.clipboard.copyFailed");
	});

	it("Should catch a synchronous clipboard failure", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.clipboard.setText = () => {
			throw new Error("runtime unavailable");
		};

		// Act
		const copied = await copyTextToClipboard("selected output");

		// Assert
		expect(copied).toBe(false);
		expect(consoleError).toHaveBeenCalledWith(
			"Failed to copy text to the clipboard",
			expect.any(Error),
		);
		expect(toastError).toHaveBeenCalledWith(
			"toast.clipboard.copyFailed: runtime unavailable",
		);
	});

	it("Should catch an asynchronous clipboard failure", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.clipboard.setText = async () => {
			throw new Error("clipboard locked");
		};

		// Act
		const copied = await copyTextToClipboard("selected output");

		// Assert
		expect(copied).toBe(false);
		expect(toastError).toHaveBeenCalledWith(
			"toast.clipboard.copyFailed: clipboard locked",
		);
	});
});
