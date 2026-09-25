import { screen } from "@testing-library/react";
import { useState } from "react";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { EditCommandGroupModal } from "@/components/modals/CommandGroup/EditCommandGroupModal.tsx";
import { resetBackendServices } from "@/contracts/service.ts";
import type { CommandGroup } from "@/contracts/types.ts";
import { commandStore } from "@/store/commandStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import { resetStores } from "@/testing/stores.ts";

const Host = ({ commandGroup }: { commandGroup: CommandGroup }) => {
	const [open, setOpen] = useState(true);

	return (
		<EditCommandGroupModal
			commandGroup={commandGroup}
			open={open}
			setOpen={setOpen}
		/>
	);
};

describe("EditCommandGroupModal", () => {
	const build = new CommandBuilder().withId("cmd-1").withName("Build").build();
	const test = new CommandBuilder().withId("cmd-2").withName("Test").build();
	const group = new CommandGroupBuilder()
		.withId("group-1")
		.withName("CI")
		.withCommands(build)
		.build();

	beforeEach(async () => {
		await installTranslations();
		resetStores();
		commandStore.setState({ commands: [build, test] });
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should show the group's name in the form", () => {
		// Arrange
		installInMemoryBackend({ commands: [build, test], commandGroups: [group] });

		// Act
		renderWithProviders(<Host commandGroup={group} />);

		// Assert
		expect(screen.getByLabelText("commandGroupForm.nameLabel")).toHaveValue(
			"CI",
		);
	});

	it("Should save the edited name and the added command and close", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commands: [build, test],
			commandGroups: [group],
		});
		const { user } = renderWithProviders(<Host commandGroup={group} />);
		const nameInput = screen.getByLabelText("commandGroupForm.nameLabel");

		// Act
		await user.clear(nameInput);
		await user.type(nameInput, "Checks");
		await user.click(
			screen.getByRole("button", { name: "commandGroupForm.addCommand" }),
		);
		await user.click(screen.getByRole("button", { name: "common.save" }));

		// Assert
		expect(backend.state.commandGroups).toEqual([
			{ ...group, name: "Checks", commandIds: ["cmd-1", "cmd-2"] },
		]);
		expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
	});

	it("Should stay open with the typed input when the backend rejects the edit", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commands: [build, test],
			commandGroups: [group],
		});
		backend.data.editCommandGroup = async () => {
			throw new Error("boom");
		};
		const { user } = renderWithProviders(<Host commandGroup={group} />);
		const nameInput = screen.getByLabelText("commandGroupForm.nameLabel");

		// Act
		await user.clear(nameInput);
		await user.type(nameInput, "Checks");
		await user.click(screen.getByRole("button", { name: "common.save" }));
		// Assert
		expect(screen.getByRole("dialog")).toBeInTheDocument();
		expect(screen.getByLabelText("commandGroupForm.nameLabel")).toHaveValue(
			"Checks",
		);
		expect(backend.state.commandGroups).toEqual([group]);
	});

	it("Should show the validation error and save nothing when every command is removed", async () => {
		// Arrange
		const backend = installInMemoryBackend({
			commands: [build, test],
			commandGroups: [group],
		});
		const { user } = renderWithProviders(<Host commandGroup={group} />);

		// Act
		await user.click(
			screen.getByRole("button", { name: "commandGroupForm.removeCommand" }),
		);
		await user.click(screen.getByRole("button", { name: "common.save" }));

		// Assert
		expect(
			screen.getByText("commandGroupForm.validation.commandsRequired"),
		).toBeInTheDocument();
		expect(backend.state.commandGroups).toEqual([group]);
	});
});
