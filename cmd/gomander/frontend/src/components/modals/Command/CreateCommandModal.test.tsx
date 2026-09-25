import { screen } from "@testing-library/react";
import { useState } from "react";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { CreateCommandModal } from "@/components/modals/Command/CreateCommandModal.tsx";
import { resetBackendServices } from "@/contracts/service.ts";
import { projectStore } from "@/store/projectStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { ProjectBuilder } from "@/testing/builders/project.ts";
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
			<CreateCommandModal open={open} setOpen={setOpen} />
		</>
	);
};

describe("CreateCommandModal", () => {
	beforeEach(async () => {
		await installTranslations();
		resetStores();
		projectStore.setState({
			projectInfo: new ProjectBuilder().withId("project-1").build(),
		});
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should create the command in the opened project, close and reset the form", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		const { user } = renderWithProviders(<Host />);

		// Act
		await user.type(screen.getByLabelText("commandForm.nameLabel"), "Build");
		await user.type(
			screen.getByLabelText("commandForm.commandLabel"),
			"make build",
		);
		await user.click(screen.getByRole("button", { name: "common.create" }));

		// Assert
		expect(backend.state.commands).toEqual([
			expect.objectContaining({
				projectId: "project-1",
				name: "Build",
				command: "make build",
			}),
		]);
		expect(screen.queryByRole("dialog")).not.toBeInTheDocument();

		await user.click(screen.getByRole("button", { name: "reopen" }));
		expect(screen.getByLabelText("commandForm.nameLabel")).toHaveValue("");
		expect(screen.getByLabelText("commandForm.commandLabel")).toHaveValue("");
	});

	it("Should stay open with the typed input when the backend rejects the command", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.addCommand = async () => {
			throw new Error("boom");
		};
		const { user } = renderWithProviders(<Host />);

		// Act
		await user.type(screen.getByLabelText("commandForm.nameLabel"), "Build");
		await user.type(
			screen.getByLabelText("commandForm.commandLabel"),
			"make build",
		);
		await user.click(screen.getByRole("button", { name: "common.create" }));
		// Assert
		expect(screen.getByRole("dialog")).toBeInTheDocument();
		expect(screen.getByLabelText("commandForm.nameLabel")).toHaveValue("Build");
		expect(screen.getByLabelText("commandForm.commandLabel")).toHaveValue(
			"make build",
		);
	});

	it("Should show the validation errors and create nothing when the required fields are empty", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		const { user } = renderWithProviders(<Host />);

		// Act
		await user.click(screen.getByRole("button", { name: "common.create" }));

		// Assert
		expect(
			screen.getByText("commandForm.validation.nameRequired"),
		).toBeInTheDocument();
		expect(
			screen.getByText("commandForm.validation.commandRequired"),
		).toBeInTheDocument();
		expect(backend.state.commands).toEqual([]);
		expect(screen.getByRole("dialog")).toBeInTheDocument();
	});
});
