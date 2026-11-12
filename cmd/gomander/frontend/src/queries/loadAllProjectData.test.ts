import { afterEach, beforeEach, describe, expect, it } from "vitest";

import type { InMemoryBackend } from "@/contracts/adapters/inMemory.ts";
import { resetBackendServices } from "@/contracts/service.ts";
import { loadAllProjectData } from "@/queries/loadAllProjectData.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { projectStore } from "@/store/projectStore.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";
import { CommandGroupBuilder } from "@/testing/builders/commandGroup.ts";
import { ProjectBuilder } from "@/testing/builders/project.ts";
import { resetStores } from "@/testing/stores.ts";

describe("loadAllProjectData", () => {
	const sut = loadAllProjectData;
	let backend: InMemoryBackend;

	beforeEach(() => {
		resetStores();
		backend = installInMemoryBackend({
			currentProject: new ProjectBuilder().withId("p-1").build(),
		});
	});

	afterEach(() => {
		resetBackendServices();
	});

	it("Should load the commands, the groups and then the opened project", async () => {
		// Arrange
		backend.state.commands = [new CommandBuilder().withId("cmd-1").build()];
		backend.state.commandGroups = [
			new CommandGroupBuilder().withId("g-1").build(),
		];

		// Act
		await sut();

		// Assert
		expect(commandGroupStore.getState().commandGroups).toHaveLength(1);
		expect(projectStore.getState().projectInfo?.id).toBe("p-1");
	});

	it("Should reject without loading the project when a query fails", async () => {
		// Arrange
		backend.data.getCommands = async () => {
			throw new Error("boom");
		};

		// Act
		const loading = sut();

		// Assert
		await expect(loading).rejects.toThrow("boom");
		expect(projectStore.getState().projectInfo).toBeNull();
	});

	it("Should not settle while the other query is still running", async () => {
		// Arrange
		backend.data.getCommands = async () => {
			throw new Error("boom");
		};
		let finishGroups: () => void = () => {};
		backend.data.getCommandGroups = () =>
			new Promise((resolve) => {
				finishGroups = () => resolve([]);
			});
		let settled = false;

		// Act
		const loading = sut()
			.catch(() => {})
			.finally(() => {
				settled = true;
			});
		await new Promise((resolve) => setTimeout(resolve, 0));

		// Assert
		expect(settled).toBe(false);
		finishGroups();
		await loading;
		expect(settled).toBe(true);
	});
});
