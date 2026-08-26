import type { ReactElement } from "react";
import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { resetBackendServices } from "@/contracts/service.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { exportProject } from "@/useCases/project/exportProject.tsx";

describe("exportProject", () => {
	const sut = exportProject;

	const toastSuccess = vi.spyOn(toast, "success");
	const toastError = vi.spyOn(toast, "error");

	beforeEach(async () => {
		vi.clearAllMocks();
		await installTranslations();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should offer the exported file to the user under the id the toast is dismissed by", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.exportProject = async () => "/exports/a-project.json";

		// Act
		const succeeded = await sut("project-1");

		// Assert
		expect(succeeded).toBe(true);
		const [content, options] = toastSuccess.mock.calls[0] as [
			ReactElement<{ exportFilePath: string; toastId: string }>,
			{ id: string },
		];
		expect(content.props.exportFilePath).toBe("/exports/a-project.json");
		expect(content.props.toastId).toBe(options.id);
	});

	it("Should notify the user when the backend rejects the export", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.exportProject = async () => {
			throw new Error("boom");
		};

		// Act
		const succeeded = await sut("project-1");

		// Assert
		expect(succeeded).toBe(false);
		expect(toastError).toHaveBeenCalledWith("toast.project.exportFailed: boom");
		expect(toastSuccess).not.toHaveBeenCalled();
	});
});
