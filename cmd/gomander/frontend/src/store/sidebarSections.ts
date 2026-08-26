import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

export const ALL_COMMANDS_SECTION_ID = "all-commands";

const STORAGE_KEY = "sidebar-sections-open";

type OpenSections = Record<string, boolean>;

const readStoredSections = (): OpenSections => {
	try {
		const stored = window.localStorage.getItem(STORAGE_KEY);
		const parsed: unknown = stored ? JSON.parse(stored) : null;

		// `JSON.parse("null")` parses fine and would then throw on every read, so
		// the shape is checked rather than the parse alone.
		return typeof parsed === "object" && parsed !== null
			? (parsed as OpenSections)
			: {};
	} catch {
		return {};
	}
};

type SidebarSectionsStore = {
	openSections: OpenSections;
};

const sidebarSectionsStore = createStore<SidebarSectionsStore>()(() => ({
	openSections: readStoredSections(),
}));

const saveOpenSections = (openSections: OpenSections) => {
	try {
		window.localStorage.setItem(STORAGE_KEY, JSON.stringify(openSections));
	} catch {
		// Ignore write errors
	}
	sidebarSectionsStore.setState({ openSections });
};

export const isSidebarSectionOpen = (sectionId: string): boolean =>
	sidebarSectionsStore.getState().openSections[sectionId] ?? false;

export const setSidebarSectionOpen = (
	sectionId: string,
	open: boolean,
): void => {
	saveOpenSections({
		...sidebarSectionsStore.getState().openSections,
		[sectionId]: open,
	});
};

export const forgetSidebarSection = (sectionId: string): void => {
	const { [sectionId]: _forgotten, ...remaining } =
		sidebarSectionsStore.getState().openSections;

	saveOpenSections(remaining);
};

// The React binding over the plain pair above, shaped like `useState` so a call
// site names its section once.
export const useSidebarSection = (
	sectionId: string,
): readonly [boolean, (open: boolean) => void] => {
	const isOpen = useStore(
		sidebarSectionsStore,
		(state) => state.openSections[sectionId] ?? false,
	);

	return [isOpen, (open: boolean) => setSidebarSectionOpen(sectionId, open)];
};
