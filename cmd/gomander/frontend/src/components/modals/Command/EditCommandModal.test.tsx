import { screen } from "@testing-library/react";
import { useState } from "react";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { EditCommandModal } from "@/components/modals/Command/EditCommandModal.tsx";
import { resetBackendServices } from "@/contracts/service.ts";
import type { Command } from "@/contracts/types.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { installTranslations } from "@/testing/i18n.ts";
import { renderWithProviders } from "@/testing/render.tsx";
import { resetStores } from "@/testing/stores.ts";

const Host = ({ command }: { command: Command }) => {
	const [open, setOpen] = useState(true);

	return <EditCommandModal command={command} open={open} setOpen={setOpen} />;
};

describe("EditCommandModal", () => {
	const command = new CommandBuilder()
		.withId("cmd-1")
		.withName("Build")
		.withCommand("make build")
		.build();

	beforeEach(async () => {
		await installTranslations();
		resetStores();
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should show the command in the form", () => {
		// Arrange
		installInMemoryBackend({ commands: [command] });

		// Act
		renderWithProviders(<Host command={command} />);

		// Assert
		expect(screen.getByLabelText("commandForm.nameLabel")).toHaveValue("Build");
		expect(screen.getByLabelText("commandForm.commandLabel")).toHaveValue(
			"make build",
		);
	});

	it("Should save the edited command and close", async () => {
		// Arrange
		const backend = installInMemoryBackend({ commands: [command] });
		const { user } = renderWithProviders(<Host command={command} />);
		const nameInput = screen.getByLabelText("commandForm.nameLabel");

		// Act
		await user.clear(nameInput);
		await user.type(nameInput, "Compile");
		await user.click(screen.getByRole("button", { name: "common.save" }));

		// Assert
		expect(backend.state.commands).toEqual([{ ...command, name: "Compile" }]);
		expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
	});

	it("Should stay open with the typed input when the backend rejects the edit", async () => {
		// Arrange
		const backend = installInMemoryBackend({ commands: [command] });
		backend.data.editCommand = async () => {
			throw new Error("boom");
		};
		const { user } = renderWithProviders(<Host command={command} />);
		const nameInput = screen.getByLabelText("commandForm.nameLabel");

		// Act
		await user.clear(nameInput);
		await user.type(nameInput, "Compile");
		await user.click(screen.getByRole("button", { name: "common.save" }));
		// Assert
		expect(screen.getByRole("dialog")).toBeInTheDocument();
		expect(screen.getByLabelText("commandForm.nameLabel")).toHaveValue(
			"Compile",
		);
		expect(backend.state.commands).toEqual([command]);
	});

	it("Should show the validation error and save nothing when the name is cleared", async () => {
		// Arrange
		const backend = installInMemoryBackend({ commands: [command] });
		const { user } = renderWithProviders(<Host command={command} />);

		// Act
		await user.clear(screen.getByLabelText("commandForm.nameLabel"));
		await user.click(screen.getByRole("button", { name: "common.save" }));

		// Assert
		expect(
			screen.getByText("commandForm.validation.nameRequired"),
		).toBeInTheDocument();
		expect(backend.state.commands).toEqual([command]);
		expect(screen.getByRole("dialog")).toBeInTheDocument();
	});
});
