import { screen } from "@testing-library/react";
import { useState } from "react";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { ImportProjectModal } from "@/components/modals/Project/ImportProjectModal.tsx";
import { resetBackendServices } from "@/contracts/service.ts";
import type { ProjectBlueprint } from "@/contracts/types.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { ProjectBlueprintBuilder } from "@/testing/builders/projectBlueprint.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import { resetStores } from "@/testing/stores.ts";

const Host = ({ project }: { project: ProjectBlueprint }) => {
	const [open, setOpen] = useState(true);

	return (
		<ImportProjectModal
			open={open}
			onClose={() => setOpen(false)}
			project={project}
		/>
	);
};

describe("ImportProjectModal", () => {
	const build = {
		id: "cmd-1",
		name: "Build",
		command: "make build",
		workingDirectory: "",
	};
	const test = {
		id: "cmd-2",
		name: "Test",
		command: "make test",
		workingDirectory: "",
	};
	const ci = { id: "group-1", name: "CI", commandIds: ["cmd-1"] };
	const all = { id: "group-2", name: "All", commandIds: ["cmd-1", "cmd-2"] };
	const blueprint = new ProjectBlueprintBuilder()
		.withName("Gomander")
		.withWorkingDirectory("/code/gomander")
		.withCommands(build, test)
		.withCommandGroups(ci, all)
		.build();

	beforeEach(async () => {
		await installTranslations();
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should import the blueprint under the chosen name and directory and close", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		const { user } = renderWithProviders(<Host project={blueprint} />);
		const nameInput = screen.getByLabelText("projectForm.nameLabel");
		const dirInput = screen.getByLabelText("projectForm.baseDirLabel");

		// Act
		await user.clear(nameInput);
		await user.type(nameInput, "Imported");
		await user.clear(dirInput);
		await user.type(dirInput, "/elsewhere");
		await user.click(screen.getByRole("button", { name: "common.save" }));

		// Assert
		expect(backend.state.importedProjects).toEqual([blueprint]);
		expect(backend.state.projects).toEqual([
			expect.objectContaining({
				name: "Imported",
				workingDirectory: "/elsewhere",
			}),
		]);
		expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
	});

	it("Should uncheck and disable a group once all of its commands are unchecked, and import only what stays checked", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		const { user } = renderWithProviders(<Host project={blueprint} />);
		await user.click(
			screen.getByRole("button", {
				name: "modal.importProject.advancedTrigger",
			}),
		);

		// Act
		await user.click(screen.getByRole("checkbox", { name: "Build" }));

		// Assert
		const ciCheckbox = screen.getByRole("checkbox", { name: "CI" });
		expect(ciCheckbox).not.toBeChecked();
		expect(ciCheckbox).toBeDisabled();
		expect(screen.getByRole("checkbox", { name: "All" })).toBeChecked();

		await user.click(screen.getByRole("button", { name: "common.save" }));
		expect(backend.state.importedProjects).toEqual([
			{ ...blueprint, commands: [test], commandGroups: [all] },
		]);
	});

	it("Should stay open with the typed input when the backend rejects the import", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		backend.data.importProject = async () => {
			throw new Error("boom");
		};
		const { user } = renderWithProviders(<Host project={blueprint} />);
		const nameInput = screen.getByLabelText("projectForm.nameLabel");

		// Act
		await user.clear(nameInput);
		await user.type(nameInput, "Imported");
		await user.click(screen.getByRole("button", { name: "common.save" }));
		// Assert
		expect(screen.getByRole("dialog")).toBeInTheDocument();
		expect(screen.getByLabelText("projectForm.nameLabel")).toHaveValue(
			"Imported",
		);
	});

	it("Should show the validation error and import nothing when the name is cleared", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		const { user } = renderWithProviders(<Host project={blueprint} />);

		// Act
		await user.clear(screen.getByLabelText("projectForm.nameLabel"));
		await user.click(screen.getByRole("button", { name: "common.save" }));

		// Assert
		expect(
			screen.getByText("projectForm.validation.nameRequired"),
		).toBeInTheDocument();
		expect(backend.state.importedProjects).toEqual([]);
	});
});
