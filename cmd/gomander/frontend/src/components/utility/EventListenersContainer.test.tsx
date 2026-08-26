import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { EventListenersContainer } from "@/components/utility/EventListenersContainer.tsx";
import type { InMemoryBackend } from "@/contracts/adapters/inMemory.ts";
import { resetBackendServices } from "@/contracts/service.ts";
import { Event } from "@/contracts/types.ts";
import {
	forgetSidebarSection,
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
	let roots: Root[] = [];
	// The store is imported once for the whole file, so what a test opens has to
	// be forgotten afterwards rather than reloaded away. Recording it here keeps
	// that automatic for sections a later test invents.
	let openedSectionIds: string[] = [];

	const openSection = (sectionId: string) => {
		openedSectionIds.push(sectionId);
		setSidebarSectionOpen(sectionId, true);
	};

	const render = async () => {
		const root = createRoot(document.createElement("div"));
		roots.push(root);
		await act(async () => {
			root.render(<EventListenersContainer />);
		});
	};

	beforeEach(() => {
		// react-dom needs this to accept the act() calls that drive the render.
		Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true });
		installRecordingTerminals();
		backend = installInMemoryBackend();
	});

	afterEach(() => {
		act(() => {
			for (const root of roots) {
				root.unmount();
			}
		});
		roots = [];

		for (const sectionId of openedSectionIds) {
			forgetSidebarSection(sectionId);
		}
		openedSectionIds = [];
		resetTerminals();
		resetBackendServices();
	});

	it("Should forget a deleted group's sidebar section", async () => {
		// Arrange
		openSection("group-1");
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
		openSection("group-1");
		openSection("group-2");
		await render();

		// Act
		await act(async () => {
			backend.emit(Event.COMMAND_GROUP_DELETED, "group-1");
		});

		// Assert
		expect(isSidebarSectionOpen("group-2")).toBe(true);
	});
});
