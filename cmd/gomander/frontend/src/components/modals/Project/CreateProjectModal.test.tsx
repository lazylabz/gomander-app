import { screen } from "@testing-library/react";
import { useState } from "react";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { CreateProjectModal } from "@/components/modals/Project/CreateProjectModal.tsx";
import { resetBackendServices } from "@/contracts/service.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import { resetStores } from "@/testing/stores.ts";

const Host = () => {
	const [open, setOpen] = useState(true);

	return (
		<>
			<button type="button" onClick={() => setOpen(true)}>
				reopen
			</button>
			<CreateProjectModal open={open} setOpen={setOpen} />
		</>
	);
};

describe("CreateProjectModal", () => {
	beforeEach(async () => {
		await installTranslations();
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should create the project, close and reset the form", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		const { user } = renderWithProviders(<Host />);

		// Act
		await user.type(screen.getByLabelText("projectForm.nameLabel"), "Gomander");
		await user.type(
			screen.getByLabelText("projectForm.baseDirLabel"),
			"/code/gomander",
		);
		await user.click(screen.getByRole("button", { name: "common.save" }));

		// Assert
		expect(backend.state.projects).toEqual([
			expect.objectContaining({
				name: "Gomander",
				workingDirectory: "/code/gomander",
			}),
		]);
		expect(screen.queryByRole("dialog")).not.toBeInTheDocument();

		await user.click(screen.getByRole("button", { name: "reopen" }));
		expect(screen.getByLabelText("projectForm.nameLabel")).toHaveValue("");
		expect(screen.getByLabelText("projectForm.baseDirLabel")).toHaveValue("");
	});

	it("Should stay open with the typed input when the backend rejects the project", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.createProject = async () => {
			throw new Error("boom");
		};
		const { user } = renderWithProviders(<Host />);

		// Act
		await user.type(screen.getByLabelText("projectForm.nameLabel"), "Gomander");
		await user.type(
			screen.getByLabelText("projectForm.baseDirLabel"),
			"/code/gomander",
		);
		await user.click(screen.getByRole("button", { name: "common.save" }));
		// Assert
		expect(screen.getByRole("dialog")).toBeInTheDocument();
		expect(screen.getByLabelText("projectForm.nameLabel")).toHaveValue(
			"Gomander",
		);
		expect(screen.getByLabelText("projectForm.baseDirLabel")).toHaveValue(
			"/code/gomander",
		);
	});

	it("Should show the validation errors and create nothing when the form is empty", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		const { user } = renderWithProviders(<Host />);

		// Act
		await user.click(screen.getByRole("button", { name: "common.save" }));

		// Assert
		expect(
			screen.getByText("projectForm.validation.nameRequired"),
		).toBeInTheDocument();
		expect(
			screen.getByText("projectForm.validation.baseDirRequired"),
		).toBeInTheDocument();
		expect(backend.state.projects).toEqual([]);
	});
});
