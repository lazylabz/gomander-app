import { act } from "react";
import { createRoot } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { EventListenersContainer } from "@/components/utility/EventListenersContainer.tsx";
import type { InMemoryBackend } from "@/contracts/adapters/inMemory.ts";
import { resetBackendServices } from "@/contracts/service.ts";
import { Event } from "@/contracts/types.ts";
import {
	isSidebarSectionOpen,
	setSidebarSectionOpen,
} from "@/store/sidebarSections.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import {
	installRecordingTerminals,
	resetTerminals,
} from "@/testing/terminals.ts";

describe("EventListenersContainer", () => {
	let backend: InMemoryBackend;

	const render = async () => {
		const root = createRoot(document.createElement("div"));
		await act(async () => {
			root.render(<EventListenersContainer />);
		});
	};

	beforeEach(() => {
		// react-dom needs this to accept the act() calls that drive the render.
		Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true });
		localStorage.clear();
		installRecordingTerminals();
		backend = installInMemoryBackend();
	});

	afterEach(() => {
		resetTerminals();
		resetBackendServices();
	});

	it("Should forget a deleted group's sidebar section", async () => {
		// Arrange
		setSidebarSectionOpen("group-1", true);
		await render();

		// Act
		await act(async () => {
			backend.emit(Event.COMMAND_GROUP_DELETED, "group-1");
		});

		// Assert
		expect(isSidebarSectionOpen("group-1")).toBe(false);
	});

	it("Should leave the other groups' sidebar sections alone", async () => {
		// Arrange
		setSidebarSectionOpen("group-1", true);
		setSidebarSectionOpen("group-2", true);
		await render();

		// Act
		await act(async () => {
			backend.emit(Event.COMMAND_GROUP_DELETED, "group-1");
		});

		// Assert
		expect(isSidebarSectionOpen("group-2")).toBe(true);
	});
});
