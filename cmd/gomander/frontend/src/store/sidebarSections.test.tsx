import { act } from "react";
import { createRoot } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

type SidebarSections = typeof import("@/store/sidebarSections.ts");

// A fresh import re-seeds the module from storage, the way a reload does.
const load = async (): Promise<SidebarSections> => {
	vi.resetModules();
	return await import("@/store/sidebarSections.ts");
};

describe("sidebarSections", () => {
	let sut: SidebarSections;

	beforeEach(async () => {
		// react-dom needs this to accept the act() calls that drive the hook.
		Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true });
		localStorage.clear();
		sut = await load();
	});

	// A storage method left mocked by a failing test would cascade into the rest.
	afterEach(() => {
		vi.restoreAllMocks();
	});

	it("Should report a section as closed when nothing was ever stored", () => {
		// Assert
		expect(sut.isSidebarSectionOpen("group-1")).toBe(false);
	});

	it("Should report a section as open once it is set open", () => {
		// Act
		sut.setSidebarSectionOpen("group-1", true);

		// Assert
		expect(sut.isSidebarSectionOpen("group-1")).toBe(true);
	});

	it("Should keep each section's state to itself", () => {
		// Act
		sut.setSidebarSectionOpen("group-1", true);

		// Assert
		expect(sut.isSidebarSectionOpen("group-2")).toBe(false);
	});

	it("Should remember a section's state across a reload", async () => {
		// Arrange
		sut.setSidebarSectionOpen("group-1", true);

		// Act
		const reloaded = await load();

		// Assert
		expect(reloaded.isSidebarSectionOpen("group-1")).toBe(true);
	});

	it("Should leave a forgotten section closed after a reload", async () => {
		// Arrange
		sut.setSidebarSectionOpen("group-1", true);

		// Act
		sut.forgetSidebarSection("group-1");
		const reloaded = await load();

		// Assert
		expect(reloaded.isSidebarSectionOpen("group-1")).toBe(false);
	});

	// "null" parses without throwing, so it reaches the read rather than the catch.
	it.each([
		"not json",
		"null",
	])("Should fall back to closed when the stored state reads back as %s", async (corrupted) => {
		// Arrange - whatever the module wrote, replaced by the corrupted value
		sut.setSidebarSectionOpen("group-1", true);
		for (let i = 0; i < localStorage.length; i++) {
			const key = localStorage.key(i);
			if (key) {
				localStorage.setItem(key, corrupted);
			}
		}

		// Act
		const reloaded = await load();

		// Assert
		expect(reloaded.isSidebarSectionOpen("group-1")).toBe(false);
	});

	it("Should stay usable when storage refuses the write", () => {
		// Arrange
		vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {
			throw new Error("quota exceeded");
		});

		// Act
		sut.setSidebarSectionOpen("group-1", true);

		// Assert
		expect(sut.isSidebarSectionOpen("group-1")).toBe(true);
	});

	it("Should re-render a component whose section is opened elsewhere", async () => {
		// Arrange
		let renderedIsOpen: boolean | undefined;
		const Probe = () => {
			[renderedIsOpen] = sut.useSidebarSection("group-1");
			return null;
		};
		const root = createRoot(document.createElement("div"));
		await act(async () => {
			root.render(<Probe />);
		});

		// Act
		await act(async () => {
			sut.setSidebarSectionOpen("group-1", true);
		});

		// Assert
		expect(renderedIsOpen).toBe(true);

		act(() => root.unmount());
	});

	it("Should open a section through the setter the hook hands back", async () => {
		// Arrange
		let setRenderedOpen: (open: boolean) => void = () => {};
		const Probe = () => {
			[, setRenderedOpen] = sut.useSidebarSection("group-1");
			return null;
		};
		const root = createRoot(document.createElement("div"));
		await act(async () => {
			root.render(<Probe />);
		});

		// Act
		await act(async () => {
			setRenderedOpen(true);
		});

		// Assert
		expect(sut.isSidebarSectionOpen("group-1")).toBe(true);

		act(() => root.unmount());
	});
});
