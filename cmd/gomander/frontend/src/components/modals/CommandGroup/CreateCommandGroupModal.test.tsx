import { screen } from "@testing-library/react";
import { useState } from "react";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { CreateCommandGroupModal } from "@/components/modals/CommandGroup/CreateCommandGroupModal.tsx";
import { resetBackendServices } from "@/contracts/service.ts";
import { commandStore } from "@/store/commandStore.ts";
import { projectStore } from "@/store/projectStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
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
			<CreateCommandGroupModal open={open} setOpen={setOpen} />
		</>
	);
};

describe("CreateCommandGroupModal", () => {
	const build = new CommandBuilder().withId("cmd-1").withName("Build").build();
	const addCommandButton = () =>
		screen.getByRole("button", { name: "commandGroupForm.addCommand" });

	beforeEach(async () => {
		await installTranslations();
		resetStores();
		projectStore.setState({
			projectInfo: new ProjectBuilder().withId("project-1").build(),
		});
		commandStore.setState({ commands: [build] });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should create the group with the picked command, close and reset the form", async () => {
		// Arrange
		const backend = installInMemoryBackend({ commands: [build] });
		const { user } = renderWithProviders(<Host />);

		// Act
		await user.type(screen.getByLabelText("commandGroupForm.nameLabel"), "CI");
		await user.click(addCommandButton());
		await user.click(screen.getByRole("button", { name: "common.create" }));

		// Assert
		expect(backend.state.commandGroups).toEqual([
			expect.objectContaining({
				projectId: "project-1",
				name: "CI",
				commandIds: ["cmd-1"],
			}),
		]);
		expect(screen.queryByRole("dialog")).not.toBeInTheDocument();

		await user.click(screen.getByRole("button", { name: "reopen" }));
		expect(screen.getByLabelText("commandGroupForm.nameLabel")).toHaveValue("");
		expect(screen.getByText("commandGroupForm.emptyGroup")).toBeInTheDocument();
	});

	it("Should stay open with the typed input when the backend rejects the group", async () => {
		// Arrange
		const backend = installInMemoryBackend({ commands: [build] });
		backend.data.createCommandGroup = async () => {
			throw new Error("boom");
		};
		const { user } = renderWithProviders(<Host />);

		// Act
		await user.type(screen.getByLabelText("commandGroupForm.nameLabel"), "CI");
		await user.click(addCommandButton());
		await user.click(screen.getByRole("button", { name: "common.create" }));
		// Assert
		expect(screen.getByRole("dialog")).toBeInTheDocument();
		expect(screen.getByLabelText("commandGroupForm.nameLabel")).toHaveValue(
			"CI",
		);
		expect(
			screen.queryByText("commandGroupForm.emptyGroup"),
		).not.toBeInTheDocument();
	});

	it("Should show the validation errors and create nothing when the form is empty", async () => {
		// Arrange
		const backend = installInMemoryBackend({ commands: [build] });
		const { user } = renderWithProviders(<Host />);

		// Act
		await user.click(screen.getByRole("button", { name: "common.create" }));

		// Assert
		expect(
			screen.getByText("commandGroupForm.validation.nameRequired"),
		).toBeInTheDocument();
		expect(
			screen.getByText("commandGroupForm.validation.commandsRequired"),
		).toBeInTheDocument();
		expect(backend.state.commandGroups).toEqual([]);
	});
});
